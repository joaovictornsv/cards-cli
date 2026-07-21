package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestRecordDeckSession(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "portuguese"})
	if err != nil {
		t.Fatal(err)
	}

	ts := "2026-07-21T12:00:00Z"
	if err := repo.RecordDeckSession(ctx, deck.ID, ts); err != nil {
		t.Fatal(err)
	}

	row, err := repo.GetDeckStatsByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if row.SessionsCount != 1 {
		t.Fatalf("sessions_count = %d, want 1", row.SessionsCount)
	}
	if !row.LastSessionAt.Valid || row.LastSessionAt.String != ts {
		t.Fatalf("last_session_at = %v, want %q", row.LastSessionAt, ts)
	}

	ts2 := "2026-07-22T12:00:00Z"
	if err := repo.RecordDeckSessionByName(ctx, "portuguese", ts2); err != nil {
		t.Fatal(err)
	}

	row, err = repo.GetDeckStatsByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if row.SessionsCount != 2 {
		t.Fatalf("sessions_count = %d, want 2", row.SessionsCount)
	}
	if row.LastSessionAt.String != ts2 {
		t.Fatalf("last_session_at = %q, want %q", row.LastSessionAt.String, ts2)
	}
}

func TestGetDeckStatsByNameNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.GetDeckStatsByName(context.Background(), "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestRecordDeckSessionNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	err = repo.RecordDeckSession(context.Background(), 999, models.NowTimestamp())
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestDeleteDeckRemovesStats(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "temp"})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.RecordDeckSession(ctx, deck.ID, models.NowTimestamp()); err != nil {
		t.Fatal(err)
	}

	if _, err := repo.DeleteDeckByName(ctx, "temp"); err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetDeckStatsByName(ctx, "temp")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound after delete, got %v", err)
	}
}
