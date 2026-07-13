package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

var ErrNotFound = errors.New("deck not found")

var ErrDuplicateName = errors.New("deck already exists")

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
		return models.Deck{}, ErrDuplicateName
	} else if !errors.Is(err, ErrNotFound) {
		return models.Deck{}, err
	}

	res, err := r.db.sql.ExecContext(ctx,
		`INSERT INTO decks (name, created_at) VALUES (?, ?)`,
		deck.Name, deck.CreatedAt,
	)
	if err != nil {
		return models.Deck{}, fmt.Errorf("insert deck: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Deck{}, fmt.Errorf("last insert id: %w", err)
	}

	return r.GetDeckByID(ctx, id)
}

func (r *Repository) GetDeckByID(ctx context.Context, id int64) (models.Deck, error) {
	row := r.db.sql.QueryRowContext(ctx, `
		SELECT d.id, d.name, d.created_at, COUNT(c.id) AS card_count
		FROM decks d
		LEFT JOIN cards c ON c.deck_id = d.id
		WHERE d.id = ?
		GROUP BY d.id`, id)

	deck, err := scanDeck(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Deck{}, ErrNotFound
	}
	if err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func (r *Repository) GetDeckByName(ctx context.Context, name string) (models.Deck, error) {
	row := r.db.sql.QueryRowContext(ctx, `
		SELECT d.id, d.name, d.created_at, COUNT(c.id) AS card_count
		FROM decks d
		LEFT JOIN cards c ON c.deck_id = d.id
		WHERE d.name = ?
		GROUP BY d.id`, name)

	deck, err := scanDeck(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Deck{}, ErrNotFound
	}
	if err != nil {
		return models.Deck{}, err
	}
	return deck, nil
}

func (r *Repository) ListDecks(ctx context.Context) ([]models.Deck, error) {
	rows, err := r.db.sql.QueryContext(ctx, `
		SELECT d.id, d.name, d.created_at, COUNT(c.id) AS card_count
		FROM decks d
		LEFT JOIN cards c ON c.deck_id = d.id
		GROUP BY d.id
		ORDER BY d.name`)
	if err != nil {
		return nil, fmt.Errorf("list decks: %w", err)
	}
	defer rows.Close()

	var decks []models.Deck
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

func (r *Repository) DeleteDeckByName(ctx context.Context, name string) (models.Deck, error) {
	deck, err := r.GetDeckByName(ctx, name)
	if err != nil {
		return models.Deck{}, err
	}

	_, err = r.db.sql.ExecContext(ctx, `DELETE FROM decks WHERE id = ?`, deck.ID)
	if err != nil {
		return models.Deck{}, fmt.Errorf("delete deck: %w", err)
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
