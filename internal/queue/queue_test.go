package queue

import (
	"errors"
	"reflect"
	"testing"
)

func TestInsertIndex(t *testing.T) {
	opts := DefaultOptions()

	tests := []struct {
		name      string
		grade     Grade
		queueLen  int
		wantIndex int
		wantErr   error
	}{
		{name: "again at offset 2", grade: GradeAgain, queueLen: 5, wantIndex: 2},
		{name: "easy at end", grade: GradeEasy, queueLen: 4, wantIndex: 4},
		{name: "replace at end", grade: GradeReplace, queueLen: 4, wantIndex: 4},
		{name: "again clamps past end", grade: GradeAgain, queueLen: 1, wantIndex: 1},
		{name: "again on empty queue", grade: GradeAgain, queueLen: 0, wantIndex: 0},
		{name: "easy on empty queue", grade: GradeEasy, queueLen: 0, wantIndex: 0},
		{name: "invalid grade", grade: Grade("unknown"), queueLen: 3, wantErr: ErrInvalidGrade},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertIndex(tt.grade, tt.queueLen, opts)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("InsertIndex() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("InsertIndex() unexpected error: %v", err)
			}
			if got != tt.wantIndex {
				t.Fatalf("InsertIndex() = %d, want %d", got, tt.wantIndex)
			}
		})
	}
}

func TestPull(t *testing.T) {
	ids := []int64{1, 2, 3, 4, 5}

	tests := []struct {
		name          string
		limit         int
		wantBatch     []int64
		wantRemaining []int64
	}{
		{name: "pull two", limit: 2, wantBatch: []int64{1, 2}, wantRemaining: []int64{3, 4, 5}},
		{name: "pull exact length", limit: 5, wantBatch: []int64{1, 2, 3, 4, 5}, wantRemaining: nil},
		{name: "pull more than length", limit: 10, wantBatch: []int64{1, 2, 3, 4, 5}, wantRemaining: nil},
		{name: "pull zero", limit: 0, wantBatch: nil, wantRemaining: []int64{1, 2, 3, 4, 5}},
		{name: "pull negative", limit: -1, wantBatch: nil, wantRemaining: []int64{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch, remaining := Pull(ids, tt.limit)
			if !reflect.DeepEqual(batch, tt.wantBatch) {
				t.Fatalf("Pull() batch = %v, want %v", batch, tt.wantBatch)
			}
			if !reflect.DeepEqual(remaining, tt.wantRemaining) {
				t.Fatalf("Pull() remaining = %v, want %v", remaining, tt.wantRemaining)
			}
		})
	}

	t.Run("empty queue", func(t *testing.T) {
		batch, remaining := Pull(nil, 4)
		if batch != nil {
			t.Fatalf("Pull() batch = %v, want nil", batch)
		}
		if remaining != nil {
			t.Fatalf("Pull() remaining = %v, want nil", remaining)
		}
	})
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name    string
		ids     []int64
		index   int
		cardID  int64
		wantIDs []int64
	}{
		{name: "front", ids: []int64{2, 3}, index: 0, cardID: 1, wantIDs: []int64{1, 2, 3}},
		{name: "middle", ids: []int64{1, 3}, index: 1, cardID: 2, wantIDs: []int64{1, 2, 3}},
		{name: "end", ids: []int64{1, 2}, index: 2, cardID: 3, wantIDs: []int64{1, 2, 3}},
		{name: "clamp beyond end", ids: []int64{1, 2}, index: 99, cardID: 3, wantIDs: []int64{1, 2, 3}},
		{name: "clamp negative", ids: []int64{1, 2}, index: -1, cardID: 0, wantIDs: []int64{0, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Insert(tt.ids, tt.index, tt.cardID)
			if !reflect.DeepEqual(got, tt.wantIDs) {
				t.Fatalf("Insert() = %v, want %v", got, tt.wantIDs)
			}
		})
	}
}

func TestReinsertAfterGrade(t *testing.T) {
	opts := DefaultOptions()

	t.Run("invalid grade", func(t *testing.T) {
		_, err := ReinsertAfterGrade([]int64{1}, 2, Grade("bad"), opts)
		if !errors.Is(err, ErrInvalidGrade) {
			t.Fatalf("ReinsertAfterGrade() error = %v, want %v", err, ErrInvalidGrade)
		}
	})

	t.Run("again insert", func(t *testing.T) {
		got, err := ReinsertAfterGrade([]int64{1, 2, 3}, 9, GradeAgain, opts)
		if err != nil {
			t.Fatalf("ReinsertAfterGrade() unexpected error: %v", err)
		}
		want := []int64{1, 2, 9, 3}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("ReinsertAfterGrade() = %v, want %v", got, want)
		}
	})
}

func TestGoldenWalkthrough(t *testing.T) {
	labels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	idsByLabel := make(map[string]int64, len(labels))
	labelByID := make(map[int64]string, len(labels))
	for i, label := range labels {
		id := int64(i + 1)
		idsByLabel[label] = id
		labelByID[id] = label
	}

	queue := make([]int64, len(labels))
	for i, label := range labels {
		queue[i] = idsByLabel[label]
	}

	opts := DefaultOptions()
	batchSize := 4

	batch, remaining := Pull(queue, batchSize)
	if !labelsEqual(batch, labelByID, []string{"A", "B", "C", "D"}) {
		t.Fatalf("Pull() batch = %v, want [A B C D]", labelsFromIDs(batch, labelByID))
	}
	if !labelsEqual(remaining, labelByID, []string{"E", "F", "G", "H"}) {
		t.Fatalf("Pull() remaining = %v, want [E F G H]", labelsFromIDs(remaining, labelByID))
	}

	steps := []struct {
		card  string
		grade Grade
		want  []string
	}{
		{card: "A", grade: GradeEasy, want: []string{"E", "F", "G", "H", "A"}},
		{card: "B", grade: GradeAgain, want: []string{"E", "F", "B", "G", "H", "A"}},
		{card: "C", grade: GradeAgain, want: []string{"E", "F", "C", "B", "G", "H", "A"}},
		{card: "D", grade: GradeEasy, want: []string{"E", "F", "C", "B", "G", "H", "A", "D"}},
	}

	current := remaining
	for _, step := range steps {
		var err error
		current, err = ReinsertAfterGrade(current, idsByLabel[step.card], step.grade, opts)
		if err != nil {
			t.Fatalf("ReinsertAfterGrade(%s) error: %v", step.card, err)
		}
		if !labelsEqual(current, labelByID, step.want) {
			t.Fatalf("after %s %s: queue = %v, want %v",
				step.card, step.grade, labelsFromIDs(current, labelByID), step.want)
		}
	}

	nextBatch, _ := Pull(current, batchSize)
	if !labelsEqual(nextBatch, labelByID, []string{"E", "F", "C", "B"}) {
		t.Fatalf("next batch = %v, want [E F C B]", labelsFromIDs(nextBatch, labelByID))
	}
}

func labelsFromIDs(ids []int64, labelByID map[int64]string) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = labelByID[id]
	}
	return out
}

func labelsEqual(ids []int64, labelByID map[int64]string, want []string) bool {
	got := labelsFromIDs(ids, labelByID)
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
