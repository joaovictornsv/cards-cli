package db

import (
	"context"
	"errors"
	"math/rand"
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

func TestListQueueCardIDsByDeck(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	deck, cards := setupDeckWithCards(t, repo, ctx)

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{cards[2].ID, cards[1].ID, cards[0].ID}
	for i, id := range want {
		if ids[i] != id {
			t.Fatalf("ids[%d] = %d, want %d", i, ids[i], id)
		}
	}

	_, err = repo.ListQueueCardIDsByDeck(ctx, "missing")
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}

	_ = deck
}

func TestReplaceDeckQueue(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	deck, cards := setupDeckWithCards(t, repo, ctx)

	reordered := []int64{cards[0].ID, cards[2].ID, cards[1].ID}
	if err := repo.ReplaceDeckQueue(ctx, deck.ID, reordered); err != nil {
		t.Fatal(err)
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	for i, id := range reordered {
		if ids[i] != id {
			t.Fatalf("ids[%d] = %d, want %d", i, ids[i], id)
		}
	}

	entries, err := repo.ListQueueByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if entries[0].FrontPreview != "first" || entries[1].FrontPreview != "third" || entries[2].FrontPreview != "second" {
		t.Fatalf("unexpected queue order: %+v", entries)
	}
}

func TestReplaceDeckQueueEmpty(t *testing.T) {
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

	if err := repo.ReplaceDeckQueue(ctx, deck.ID, nil); err != nil {
		t.Fatal(err)
	}

	ids, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 0 {
		t.Fatalf("expected empty queue, got %v", ids)
	}
}

func TestShuffleDeckQueueNoop(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(42))

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "empty"}); err != nil {
		t.Fatal(err)
	}

	result, err := repo.ShuffleDeckQueue(ctx, "empty", rng)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "noop" || result.CardCount != 0 {
		t.Fatalf("expected noop with 0 cards, got %+v", result)
	}

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "single"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.CreateCard(ctx, deck.Name, models.Card{Front: "one", Back: "uno"}); err != nil {
		t.Fatal(err)
	}

	result, err = repo.ShuffleDeckQueue(ctx, "single", rng)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "noop" || result.CardCount != 1 {
		t.Fatalf("expected noop with 1 card, got %+v", result)
	}
}

func TestShuffleDeckQueueChangesOrder(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	before, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(42))
	result, err := repo.ShuffleDeckQueue(ctx, "portuguese", rng)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "shuffled" || result.CardCount != 3 {
		t.Fatalf("expected shuffled with 3 cards, got %+v", result)
	}

	after, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}

	if slicesEqual(before, after) {
		t.Fatalf("expected order to change, got %v", after)
	}
	if !isPermutation(before, after) {
		t.Fatalf("expected permutation of %v, got %v", before, after)
	}

	_ = cards
}

func TestShuffleDeckQueueSeededDeterministic(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupDeckWithCards(t, repo, ctx)

	original, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}

	rng1 := rand.New(rand.NewSource(99))
	if _, err := repo.ShuffleDeckQueue(ctx, "portuguese", rng1); err != nil {
		t.Fatal(err)
	}
	first, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}

	deck, err := repo.GetDeckByName(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.ReplaceDeckQueue(ctx, deck.ID, original); err != nil {
		t.Fatal(err)
	}

	rng2 := rand.New(rand.NewSource(99))
	if _, err := repo.ShuffleDeckQueue(ctx, "portuguese", rng2); err != nil {
		t.Fatal(err)
	}
	second, err := repo.ListQueueCardIDsByDeck(ctx, "portuguese")
	if err != nil {
		t.Fatal(err)
	}

	if !slicesEqual(first, second) {
		t.Fatalf("expected same order with same seed, got %v vs %v", first, second)
	}
}

func TestShuffleDeckQueueNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(1))

	_, err = repo.ShuffleDeckQueue(ctx, "missing", rng)
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func slicesEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isPermutation(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[int64]int, len(a))
	for _, id := range a {
		counts[id]++
	}
	for _, id := range b {
		counts[id]--
		if counts[id] < 0 {
			return false
		}
	}
	return true
}
