package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestListQueueByDeckEmpty(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "portuguese"}); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.ListQueueByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if entries == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestListQueueByDeckAfterAdds(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	entries, err := repo.ListQueueByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	want := []struct {
		pos   int
		id    int64
		front string
	}{
		{0, cards[2].ID, "third"},
		{1, cards[1].ID, "second"},
		{2, cards[0].ID, "first"},
	}
	for i, w := range want {
		if entries[i].Position != w.pos {
			t.Fatalf("entry %d: expected position %d, got %d", i, w.pos, entries[i].Position)
		}
		if entries[i].ID != w.id {
			t.Fatalf("entry %d: expected id %d, got %d", i, w.id, entries[i].ID)
		}
		if entries[i].FrontPreview != w.front {
			t.Fatalf("entry %d: expected front %q, got %q", i, w.front, entries[i].FrontPreview)
		}
	}
}

func TestListQueueByDeckAfterDelete(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if _, err := repo.DeleteCard(ctx, "portuguese", cards[1].ID); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.ListQueueByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Position != 0 || entries[0].FrontPreview != "third" {
		t.Fatalf("expected third at position 0, got %+v", entries[0])
	}
	if entries[1].Position != 1 || entries[1].FrontPreview != "first" {
		t.Fatalf("expected first at position 1, got %+v", entries[1])
	}
}

func TestListQueueByDeckNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	_, err = repo.ListQueueByDeck(ctx, "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}
