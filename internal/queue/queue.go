package queue

import (
	"errors"
	"fmt"
)

var ErrInvalidGrade = errors.New("invalid grade")

type Grade string

const (
	GradeAgain Grade = "again"
	GradeHard  Grade = "hard"
	GradeEasy  Grade = "easy"
)

type Options struct {
	AgainOffset int
	HardOffset  int
}

func DefaultOptions() Options {
	return Options{
		AgainOffset: 2,
		HardOffset:  5,
	}
}

// InsertIndex returns the 0-based insert index for grade given the current queue
// length (after the reviewed card was removed).
func InsertIndex(grade Grade, queueLen int, opts Options) (int, error) {
	switch grade {
	case GradeEasy:
		return queueLen, nil
	case GradeAgain:
		return clampIndex(opts.AgainOffset, queueLen), nil
	case GradeHard:
		return clampIndex(opts.HardOffset, queueLen), nil
	default:
		return 0, fmt.Errorf("%w: %q", ErrInvalidGrade, grade)
	}
}

func clampIndex(index, queueLen int) int {
	if index < 0 {
		return 0
	}
	if index > queueLen {
		return queueLen
	}
	return index
}

// Pull removes up to limit card IDs from the front of ids.
func Pull(ids []int64, limit int) (batch, remaining []int64) {
	if limit <= 0 || len(ids) == 0 {
		return nil, ids
	}
	if limit >= len(ids) {
		return append([]int64(nil), ids...), nil
	}
	return append([]int64(nil), ids[:limit]...), append([]int64(nil), ids[limit:]...)
}

// Insert inserts cardID at index, clamping index to [0, len(ids)].
func Insert(ids []int64, index int, cardID int64) []int64 {
	index = clampIndex(index, len(ids))
	out := make([]int64, 0, len(ids)+1)
	out = append(out, ids[:index]...)
	out = append(out, cardID)
	out = append(out, ids[index:]...)
	return out
}

// ReinsertAfterGrade inserts cardID into ids using grade rules. The card must
// already have been removed from the queue before grading.
func ReinsertAfterGrade(ids []int64, cardID int64, grade Grade, opts Options) ([]int64, error) {
	index, err := InsertIndex(grade, len(ids), opts)
	if err != nil {
		return nil, err
	}
	return Insert(ids, index, cardID), nil
}
