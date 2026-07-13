package study

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/queue"
)

func TestSessionPersistDirect(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "two"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := repo.CreateCard(ctx, "two", models.Card{Front: "B", Back: "B"})
	if err != nil {
		t.Fatal(err)
	}
	a, err := repo.CreateCard(ctx, "two", models.Card{Front: "A", Back: "A"})
	if err != nil {
		t.Fatal(err)
	}

	sess := &Session{Store: NewDBStore(repo)}
	if err := sess.persist(ctx, deck.ID, []int64{b.ID}, []int64{a.ID}); err != nil {
		t.Fatal(err)
	}
	entries, err := repo.ListQueueByDeck(ctx, "two")
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].FrontPreview != "B" {
		t.Fatalf("expected B at front, got %+v", entries)
	}
}

func TestSessionTwoCardReorder(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "two"}); err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{"B", "A"} {
		if _, err := repo.CreateCard(ctx, "two", models.Card{Front: label, Back: label}); err != nil {
			t.Fatal(err)
		}
	}

	sess := &Session{
		DeckName: "two",
		Out:      &bytes.Buffer{},
		Store:    NewDBStore(repo),
		Input: NewScriptedInput([]queue.Grade{queue.GradeEasy}),
		Opts:     Options{BatchSize: 1, QueueOpts: queue.DefaultOptions()},
	}
	if err := sess.Run(ctx); err != nil {
		t.Fatal(err)
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "two")
	if err != nil {
		t.Fatal(err)
	}
	entries, err := repo.ListQueueByDeck(ctx, "two")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 || entries[0].FrontPreview != "B" || entries[1].FrontPreview != "A" {
		t.Fatalf("expected [B, A], got queue ids %v entries %+v", ids, entries)
	}
}

func TestReplaceDeckQueueDirect(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "two"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := repo.CreateCard(ctx, "two", models.Card{Front: "B", Back: "B"})
	if err != nil {
		t.Fatal(err)
	}
	a, err := repo.CreateCard(ctx, "two", models.Card{Front: "A", Back: "A"})
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.ReplaceDeckQueue(ctx, deck.ID, []int64{b.ID, a.ID}); err != nil {
		t.Fatal(err)
	}
	entries, err := repo.ListQueueByDeck(ctx, "two")
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].FrontPreview != "B" || entries[1].FrontPreview != "A" {
		t.Fatalf("expected [B, A], got %+v", entries)
	}
}

func TestSessionGoldenWalkthrough(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "walk"}); err != nil {
		t.Fatal(err)
	}

	labels := []string{"H", "G", "F", "E", "D", "C", "B", "A"}
	idByLabel := make(map[string]int64, len(labels))
	for _, label := range labels {
		card, err := repo.CreateCard(ctx, "walk", models.Card{
			Front: label,
			Back:  label + " back",
		})
		if err != nil {
			t.Fatal(err)
		}
		idByLabel[label] = card.ID
	}

	var out bytes.Buffer
	sess := &Session{
		DeckName: "walk",
		Out:      &out,
		Store:    NewDBStore(repo),
		Input: NewScriptedInput([]queue.Grade{
			queue.GradeEasy,
			queue.GradeAgain,
			queue.GradeHard,
			queue.GradeEasy,
		}),
		Opts: Options{
			BatchSize: 4,
			QueueOpts: queue.DefaultOptions(),
		},
	}

	if err := sess.Run(ctx); err != nil {
		t.Fatal(err)
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "walk")
	if err != nil {
		t.Fatal(err)
	}

	wantLabels := []string{"E", "F", "B", "G", "H", "C", "A", "D"}
	if len(ids) != len(wantLabels) {
		t.Fatalf("queue length = %d, want %d", len(ids), len(wantLabels))
	}
	labelByID := make(map[int64]string, len(idByLabel))
	for label, id := range idByLabel {
		labelByID[id] = label
	}
	for i, id := range ids {
		got := labelByID[id]
		if got != wantLabels[i] {
			t.Fatalf("position %d: got %q, want %q (full queue: %v)", i, got, wantLabels[i], labelsFromIDs(ids, labelByID))
		}
	}

	if !strings.Contains(out.String(), "[1/4] A") {
		t.Fatalf("expected progress output, got:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "Session complete.") {
		t.Fatalf("expected session complete message, got:\n%s", out.String())
	}
}

func TestSessionSmallDeckBatchClamp(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "small"}); err != nil {
		t.Fatal(err)
	}
	card, err := repo.CreateCard(ctx, "small", models.Card{Front: "only", Back: "one"})
	if err != nil {
		t.Fatal(err)
	}

	var out bytes.Buffer
	sess := &Session{
		DeckName: "small",
		Out:      &out,
		Store:    NewDBStore(repo),
		Input: NewScriptedInput([]queue.Grade{queue.GradeEasy}),
		Opts: Options{
			BatchSize: 4,
			QueueOpts: queue.DefaultOptions(),
		},
	}

	if err := sess.Run(ctx); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out.String(), "[1/1]") {
		t.Fatalf("expected [1/1] progress, got:\n%s", out.String())
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "small")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 || ids[0] != card.ID {
		t.Fatalf("expected single card at end, got %v", ids)
	}
}

func TestSessionDeckNotFound(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	sess := &Session{
		DeckName: "missing",
		Out:      &bytes.Buffer{},
		Store:    NewDBStore(db.NewRepository(database)),
		Input: NewScriptedInput(nil),
		Opts:     Options{BatchSize: 4, QueueOpts: queue.DefaultOptions()},
	}

	err = sess.Run(context.Background())
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestSessionQuitPersistsPendingBatch(t *testing.T) {
	database, err := db.OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "quit"}); err != nil {
		t.Fatal(err)
	}
	for _, label := range []string{"D", "C", "B", "A"} {
		if _, err := repo.CreateCard(ctx, "quit", models.Card{Front: label, Back: label}); err != nil {
			t.Fatal(err)
		}
	}

	sess := &Session{
		DeckName: "quit",
		Out:      &bytes.Buffer{},
		Store:    NewDBStore(repo),
		Input: NewScriptedInput([]queue.Grade{queue.GradeEasy}).WithQuitAt(1),
		Opts: Options{
			BatchSize: 4,
			QueueOpts: queue.DefaultOptions(),
		},
	}

	if err := sess.Run(ctx); err != nil {
		t.Fatalf("expected nil after quit persist, got %v", err)
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "quit")
	if err != nil {
		t.Fatal(err)
	}
	// After grading A easy: pending was [B,C,D], tail was [A]
	// Quit before B reveal: pending [B,C,D] + tail [A]
	if len(ids) != 4 {
		t.Fatalf("expected 4 cards in queue, got %d", len(ids))
	}
}

func labelsFromIDs(ids []int64, labelByID map[int64]string) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = labelByID[id]
	}
	return out
}
