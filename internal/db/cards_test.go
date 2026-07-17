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

	cards, err := repo.ListCardsByDeck(ctx, "portuguese", false)
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

	cards, err := repo.ListCardsByDeck(ctx, "portuguese", false)
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
	_, err = repo.ListCardsByDeck(context.Background(), "missing", false)
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func setupDeckWithCards(t *testing.T, repo *Repository, ctx context.Context) (deck models.Deck, cards []models.Card) {
	t.Helper()
	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "portuguese"})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ front, back string }{
		{"first", "first back"},
		{"second", "second back"},
		{"third", "third back"},
	} {
		card, err := repo.CreateCard(ctx, "portuguese", models.Card{Front: c.front, Back: c.back})
		if err != nil {
			t.Fatal(err)
		}
		cards = append(cards, card)
	}
	return deck, cards
}

func TestGetCardByDeckAndID(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	got, err := repo.GetCardByDeckAndID(ctx, "portuguese", cards[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Front != "second" {
		t.Fatalf("expected front second, got %q", got.Front)
	}
	if got.Back != "second back" {
		t.Fatalf("expected back second back, got %q", got.Back)
	}
}

func TestGetCardByDeckAndIDDeckNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	_, err = repo.GetCardByDeckAndID(context.Background(), "missing", 1)
	if !errors.Is(err, ErrDeckNotFound) {
		t.Fatalf("expected ErrDeckNotFound, got %v", err)
	}
}

func TestGetCardByDeckAndIDCardNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupDeckWithCards(t, repo, ctx)

	_, err = repo.GetCardByDeckAndID(ctx, "portuguese", 9999)
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}
}

func TestGetCardByDeckAndIDWrongDeck(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if _, err := repo.CreateDeck(ctx, models.Deck{Name: "spanish"}); err != nil {
		t.Fatal(err)
	}

	_, err = repo.GetCardByDeckAndID(ctx, "spanish", cards[0].ID)
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}
}

func TestUpdateCardFrontOnly(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	front := "updated second"
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[1].ID, &front, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Front != "updated second" {
		t.Fatalf("expected updated front, got %q", updated.Front)
	}
	if updated.Back != "second back" {
		t.Fatalf("expected back unchanged, got %q", updated.Back)
	}

	reloaded, err := repo.GetCardByDeckAndID(ctx, "portuguese", cards[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Front != "updated second" {
		t.Fatalf("expected persisted front update, got %q", reloaded.Front)
	}
}

func TestUpdateCardBackOnly(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	back := "updated back"
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[0].ID, nil, &back, nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Back != "updated back" {
		t.Fatalf("expected updated back, got %q", updated.Back)
	}
}

func TestUpdateCardBoth(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	front := "new front"
	back := "new back"
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[2].ID, &front, &back, nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Front != "new front" || updated.Back != "new back" {
		t.Fatalf("unexpected update: %+v", updated)
	}
}

func TestUpdateCardValidatesEmptyFront(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	front := "   "
	_, err = repo.UpdateCard(ctx, "portuguese", cards[0].ID, &front, nil, nil)
	if !errors.Is(err, models.ErrCardFrontRequired) {
		t.Fatalf("expected ErrCardFrontRequired, got %v", err)
	}
}

func TestUpdateCardDuplicateFrontAllowed(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	front := "first"
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[1].ID, &front, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Front != "first" {
		t.Fatalf("expected duplicate front allowed, got %q", updated.Front)
	}
}

func TestDeleteCardRemovesFromQueue(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	deck, cards := setupDeckWithCards(t, repo, ctx)

	// Queue order after inserts: third(0), second(1), first(2)
	deleted, err := repo.DeleteCard(ctx, "portuguese", cards[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleted.Front != "second" {
		t.Fatalf("expected deleted card second, got %q", deleted.Front)
	}

	var count int
	if err := database.SQL().QueryRowContext(ctx, `SELECT COUNT(*) FROM cards WHERE id = ?`, cards[1].ID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected card removed, got count %d", count)
	}

	rows, err := database.SQL().QueryContext(ctx, `
		SELECT q.position, c.front
		FROM queue q
		JOIN cards c ON c.id = q.card_id
		WHERE q.deck_id = ?
		ORDER BY q.position`,
		deck.ID,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	type entry struct {
		pos   int
		front string
	}
	var queue []entry
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.pos, &e.front); err != nil {
			t.Fatal(err)
		}
		queue = append(queue, e)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if len(queue) != 2 {
		t.Fatalf("expected 2 queue entries, got %d", len(queue))
	}
	if queue[0].pos != 0 || queue[0].front != "third" {
		t.Fatalf("expected third at position 0, got %+v", queue[0])
	}
	if queue[1].pos != 1 || queue[1].front != "first" {
		t.Fatalf("expected first at position 1, got %+v", queue[1])
	}
}

func TestDeleteCardCompactsLargeQueue(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()

	deck, err := repo.CreateDeck(ctx, models.Deck{Name: "large"})
	if err != nil {
		t.Fatal(err)
	}

	const n = 70
	cards := make([]models.Card, n)
	for i := 0; i < n; i++ {
		card, err := repo.CreateCard(ctx, "large", models.Card{
			Front: "card",
			Back:  "back",
		})
		if err != nil {
			t.Fatal(err)
		}
		cards[i] = card
	}

	// Delete from the middle of a large queue (regression for UNIQUE collisions).
	middle := n / 2
	if _, err := repo.DeleteCard(ctx, "large", cards[middle].ID); err != nil {
		t.Fatal(err)
	}

	rows, err := database.SQL().QueryContext(ctx, `
		SELECT position FROM queue WHERE deck_id = ? ORDER BY position`,
		deck.ID,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	prev := -1
	count := 0
	for rows.Next() {
		var pos int
		if err := rows.Scan(&pos); err != nil {
			t.Fatal(err)
		}
		if pos != prev+1 {
			t.Fatalf("expected contiguous positions, got %d after %d", pos, prev)
		}
		prev = pos
		count++
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if count != n-1 {
		t.Fatalf("expected %d queue entries, got %d", n-1, count)
	}
}

func TestDeleteCardNotFound(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	setupDeckWithCards(t, repo, ctx)

	_, err = repo.DeleteCard(ctx, "portuguese", 9999)
	if !errors.Is(err, ErrCardNotFound) {
		t.Fatalf("expected ErrCardNotFound, got %v", err)
	}
}

func TestSetReplaceEligible(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if err := repo.SetReplaceEligible(ctx, "portuguese", cards[0].ID, true); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetCardByDeckAndID(ctx, "portuguese", cards[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	if !got.ReplaceEligible {
		t.Fatal("expected replace_eligible true")
	}
}

func TestListCardsByDeckReplaceEligibleFilter(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if err := repo.SetReplaceEligible(ctx, "portuguese", cards[1].ID, true); err != nil {
		t.Fatal(err)
	}

	flagged, err := repo.ListCardsByDeck(ctx, "portuguese", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(flagged) != 1 || flagged[0].ID != cards[1].ID {
		t.Fatalf("expected 1 flagged card id %d, got %+v", cards[1].ID, flagged)
	}
	if !flagged[0].ReplaceEligible {
		t.Fatal("expected replace_eligible in summary")
	}
}

func TestUpdateCardReplaceEligibleOnly(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if err := repo.SetReplaceEligible(ctx, "portuguese", cards[0].ID, true); err != nil {
		t.Fatal(err)
	}

	eligible := false
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[0].ID, nil, nil, &eligible)
	if err != nil {
		t.Fatal(err)
	}
	if updated.ReplaceEligible {
		t.Fatal("expected replace_eligible cleared")
	}
	if updated.Front != "first" {
		t.Fatalf("expected front unchanged, got %q", updated.Front)
	}
}

func TestUpdateCardFrontPreservesReplaceEligible(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	repo := NewRepository(database)
	ctx := context.Background()
	_, cards := setupDeckWithCards(t, repo, ctx)

	if err := repo.SetReplaceEligible(ctx, "portuguese", cards[0].ID, true); err != nil {
		t.Fatal(err)
	}

	front := "updated first"
	updated, err := repo.UpdateCard(ctx, "portuguese", cards[0].ID, &front, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !updated.ReplaceEligible {
		t.Fatal("expected replace_eligible preserved after front edit")
	}
}
