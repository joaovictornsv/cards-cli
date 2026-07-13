package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestCreateCardInsertsAtQueueFront(t *testing.T) {
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

	first, err := repo.CreateCard(ctx, "portuguese", models.Card{
		Front: "first",
		Back:  "first back",
	})
	if err != nil {
		t.Fatal(err)
	}

	second, err := repo.CreateCard(ctx, "portuguese", models.Card{
		Front: "second",
		Back:  "second back",
	})
	if err != nil {
		t.Fatal(err)
	}

	if first.ID <= 0 || second.ID <= 0 {
		t.Fatalf("expected positive card ids, got %d and %d", first.ID, second.ID)
	}
	if first.DeckID != deck.ID || second.DeckID != deck.ID {
		t.Fatalf("expected deck_id %d, got %d and %d", deck.ID, first.DeckID, second.DeckID)
	}

	var firstPos, secondPos int
	err = database.SQL().QueryRowContext(ctx, `
		SELECT position FROM queue WHERE deck_id = ? AND card_id = ?`,
		deck.ID, first.ID,
	).Scan(&firstPos)
	if err != nil {
		t.Fatal(err)
	}
	err = database.SQL().QueryRowContext(ctx, `
		SELECT position FROM queue WHERE deck_id = ? AND card_id = ?`,
		deck.ID, second.ID,
	).Scan(&secondPos)
	if err != nil {
		t.Fatal(err)
	}
	if secondPos != 0 {
		t.Fatalf("expected newest card at position 0, got %d", secondPos)
	}
	if firstPos != 1 {
		t.Fatalf("expected first card at position 1, got %d", firstPos)
	}
}

func TestCreateCardDeckNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.CreateCard(context.Background(), "missing", models.Card{
		Front: "front",
		Back:  "back",
	})
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestCreateCardValidatesBeforeDB(t *testing.T) {
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

	_, err = repo.CreateCard(ctx, "portuguese", models.Card{Front: "", Back: "back"})
	if !errors.Is(err, models.ErrCardFrontRequired) {
		t.Fatalf("expected ErrCardFrontRequired, got %v", err)
	}

	var count int
	if err := database.SQL().QueryRowContext(ctx, `SELECT COUNT(*) FROM cards`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no cards inserted, got %d", count)
	}
}

func TestListCardsByDeck(t *testing.T) {
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

	created, err := repo.CreateCard(ctx, "portuguese", models.Card{
		Front: "What is saudade?",
		Back:  "A deep emotional state of longing.",
	})
	if err != nil {
		t.Fatal(err)
	}

	cards, err := repo.ListCardsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(cards))
	}
	if cards[0].ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, cards[0].ID)
	}
	if cards[0].Front != "What is saudade?" {
		t.Fatalf("expected front %q, got %q", "What is saudade?", cards[0].Front)
	}
	if cards[0].CreatedAt == "" || cards[0].UpdatedAt == "" {
		t.Fatalf("expected timestamps, got created_at=%q updated_at=%q", cards[0].CreatedAt, cards[0].UpdatedAt)
	}
}

func TestListCardsByDeckEmpty(t *testing.T) {
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

	cards, err := repo.ListCardsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if cards == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(cards) != 0 {
		t.Fatalf("expected 0 cards, got %d", len(cards))
	}
}

func TestListCardsByDeckNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.ListCardsByDeck(context.Background(), "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}
