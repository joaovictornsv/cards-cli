package db

import (
	"context"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func setupSearchFixtures(t *testing.T, repo *Repository, ctx context.Context) {
	t.Helper()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "portuguese"}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "spanish"}); err != nil {
		t.Fatal(err)
	}

	cards := []struct {
		deck, front, back string
	}{
		{"portuguese", "What is saudade?", "A deep emotional state of longing."},
		{"portuguese", "How do you say hello?", "Olá"},
		{"spanish", "What is nostalgia?", "A sentimental longing for the past."},
	}
	for _, card := range cards {
		if _, err := repo.CreateCard(ctx, card.deck, models.Card{
			Front: card.front,
			Back:  card.back,
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSearchCardsByFront(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"saudade"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Deck != "portuguese" || results[0].Front != "What is saudade?" {
		t.Fatalf("unexpected result: %+v", results[0])
	}
}

func TestSearchCardsByBack(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"Olá"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Back != "Olá" {
		t.Fatalf("expected back Olá, got %q", results[0].Back)
	}
}

func TestSearchCardsByDeckName(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"spanish"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result from spanish deck, got %d", len(results))
	}
	if results[0].Deck != "spanish" {
		t.Fatalf("expected spanish deck, got %q", results[0].Deck)
	}
}

func TestSearchCardsORTerms(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"saudade", "nostalgia"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchCardsCaseInsensitive(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"SAUDADE"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchCardsDeckFilter(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"hello"}, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result in portuguese, got %d", len(results))
	}
	if results[0].Deck != "portuguese" {
		t.Fatalf("expected portuguese deck, got %q", results[0].Deck)
	}

	results, err = repo.SearchCards(ctx, []string{"hello"}, "spanish")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results in spanish, got %d", len(results))
	}
}

func TestSearchCardsDeckFilterNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	_, err = repo.SearchCards(ctx, []string{"hello"}, "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestSearchCardsEmptyResults(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupSearchFixtures(t, repo, ctx)

	results, err := repo.SearchCards(ctx, []string{"nonexistent"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if results == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchCardsNoTerms(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.SearchCards(context.Background(), nil, "")
	if err == nil {
		t.Fatal("expected error for missing terms")
	}
}

func TestSearchCardsEscapesLikeWildcards(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "test"}); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateCard(ctx, "test", models.Card{
		Front: "100% complete",
		Back:  "literal percent sign",
	}); err != nil {
		t.Fatal(err)
	}

	results, err := repo.SearchCards(ctx, []string{"100%"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for literal %%, got %d", len(results))
	}
}
