package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

var ErrDeckNotFound = errors.New("deck not found")

var ErrDeckDuplicateName = errors.New("deck already exists")

const deckSelectBase = `
	SELECT d.id, d.name, d.created_at, COUNT(c.id) AS card_count
	FROM decks d
	LEFT JOIN cards c ON c.deck_id = d.id`

type rowScanner interface {
	Scan(dest ...any) error
}

func (r *Repository) CreateDeck(ctx context.Context, deck models.Deck) (models.Deck, error) {
	if deck.CreatedAt == "" {
		deck.CreatedAt = models.NowTimestamp()
	}
	if err := deck.ValidateForCreate(); err != nil {
		return models.Deck{}, err
	}

	if _, err := r.GetDeckByName(ctx, deck.Name); err == nil {
		return models.Deck{}, ErrDeckDuplicateName
	} else if !errors.Is(err, ErrDeckNotFound) {
		return models.Deck{}, err
	}

	res, err := r.db.sql.ExecContext(ctx,
		`INSERT INTO decks (name, created_at) VALUES (?, ?)`,
		deck.Name, deck.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Deck{}, ErrDeckDuplicateName
		}
		return models.Deck{}, fmt.Errorf("insert deck: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Deck{}, fmt.Errorf("last insert id: %w", err)
	}

	return r.GetDeckByID(ctx, id)
}

func (r *Repository) GetDeckByID(ctx context.Context, id int64) (models.Deck, error) {
	row := r.db.sql.QueryRowContext(ctx,
		deckSelectBase+` WHERE d.id = ? GROUP BY d.id`, id)

	deck, err := scanDeck(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Deck{}, ErrDeckNotFound
	}
	if err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func (r *Repository) GetDeckByName(ctx context.Context, name string) (models.Deck, error) {
	row := r.db.sql.QueryRowContext(ctx,
		deckSelectBase+` WHERE d.name = ? GROUP BY d.id`, name)

	deck, err := scanDeck(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Deck{}, ErrDeckNotFound
	}
	if err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func (r *Repository) ListDecks(ctx context.Context) ([]models.Deck, error) {
	rows, err := r.db.sql.QueryContext(ctx,
		deckSelectBase+` GROUP BY d.id ORDER BY d.name`)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	defer rows.Close()

	decks := make([]models.Deck, 0)
	for rows.Next() {
		deck, err := scanDeck(rows)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return decks, nil
}

func (r *Repository) DeleteDeckByID(ctx context.Context, id int64) error {
	_, err := r.db.sql.ExecContext(ctx, `DELETE FROM decks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete deck: %w", err)
	}
	return nil
}

func (r *Repository) DeleteDeckByName(ctx context.Context, name string) (models.Deck, error) {
	deck, err := r.GetDeckByName(ctx, name)
	if err != nil {
		return models.Deck{}, err
	}

	if err := r.DeleteDeckByID(ctx, deck.ID); err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func scanDeck(row rowScanner) (models.Deck, error) {
	var deck models.Deck
	if err := row.Scan(&deck.ID, &deck.Name, &deck.CreatedAt, &deck.CardCount); err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
