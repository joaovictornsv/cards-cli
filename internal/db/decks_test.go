package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestRepositoryDeckCRUD(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	created, err := repo.CreateDeck(ctx, models.Deck{Name: "portuguese"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected positive id, got %d", created.ID)
	}
	if created.Name != "portuguese" {
		t.Fatalf("expected name portuguese, got %q", created.Name)
	}
	if created.CardCount != 0 {
		t.Fatalf("expected card_count 0, got %d", created.CardCount)
	}
	if created.CreatedAt == "" {
		t.Fatal("expected created_at to be set")
	}

	got, err := repo.GetDeckByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected id %d, got %d", created.ID, got.ID)
	}

	decks, err := repo.ListDecks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(decks) != 1 {
		t.Fatalf("expected 1 deck, got %d", len(decks))
	}

	_, err = repo.CreateDeck(ctx, models.Deck{Name: "portuguese"})
	if !errors.Is(err, ErrDeckDuplicateName) {
		t.Fatalf("expected ErrDeckDuplicateName, got %v", err)
	}

	deleted, err := repo.DeleteDeckByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if deleted.Name != "portuguese" {
		t.Fatalf("expected deleted name portuguese, got %q", deleted.Name)
	}

	_, err = repo.GetDeckByName(ctx, "portuguese")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestGetDeckByNameNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.GetDeckByName(context.Background(), "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestDeleteDeckByNameNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.DeleteDeckByName(context.Background(), "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestDeleteDeckCascadesCardsAndQueue(t *testing.T) {
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

	now := models.NowTimestamp()
	res, err := database.SQL().ExecContext(ctx, `
		INSERT INTO cards (deck_id, front, back, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		deck.ID, "front", "back", now, now,
	)
	if err != nil {
		t.Fatal(err)
	}
	cardID, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := database.SQL().ExecContext(ctx,
		`INSERT INTO queue (deck_id, position, card_id) VALUES (?, ?, ?)`,
		deck.ID, 0, cardID,
	); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetDeckByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if got.CardCount != 1 {
		t.Fatalf("expected card_count 1, got %d", got.CardCount)
	}

	if _, err := repo.DeleteDeckByName(ctx, "portuguese"); err != nil {
		t.Fatal(err)
	}

	for _, table := range []string{"cards", "queue"} {
		var count int
		err := database.SQL().QueryRowContext(ctx,
			`SELECT COUNT(*) FROM `+table+` WHERE deck_id = ?`, deck.ID,
		).Scan(&count)
		if err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("expected 0 rows in %s after delete, got %d", table, count)
		}
	}
}

func TestListDecksEmpty(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	decks, err := repo.ListDecks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if decks == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(decks) != 0 {
		t.Fatalf("expected 0 decks, got %d", len(decks))
	}
}

func TestIsUniqueViolation(t *testing.T) {
	if !isUniqueViolation(errors.New("UNIQUE constraint failed: decks.name")) {
		t.Fatal("expected unique violation match")
	}
	if isUniqueViolation(errors.New("other error")) {
		t.Fatal("expected no match for unrelated error")
	}
}
