package db

import (
	"context"
	"fmt"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func (r *Repository) ListQueueCardIDsByDeck(ctx context.Context, deckName string) ([]int64, error) {
	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.sql.QueryContext(ctx, `
		SELECT card_id
		FROM queue
		WHERE deck_id = ?
		ORDER BY position`,
		deck.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query queue card ids: %w", err)
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan queue card id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate queue card ids: %w", err)
	}

	return ids, nil
}

func (r *Repository) ReplaceDeckQueue(ctx context.Context, deckID int64, cardIDs []int64) error {
	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `DELETE FROM queue WHERE deck_id = ?`, deckID); err != nil {
		return fmt.Errorf("delete queue: %w", err)
	}

	for i, cardID := range cardIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO queue (deck_id, position, card_id)
			VALUES (?, ?, ?)`,
			deckID, i, cardID,
		); err != nil {
			return fmt.Errorf("insert queue row: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *Repository) ListQueueByDeck(ctx context.Context, deckName string) ([]models.QueueEntry, error) {
	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.sql.QueryContext(ctx, `
		SELECT q.position, c.id, c.front
		FROM queue q
		JOIN cards c ON c.id = q.card_id
		WHERE q.deck_id = ?
		ORDER BY q.position`,
		deck.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("query queue: %w", err)
	}
	defer rows.Close()

	entries := []models.QueueEntry{}
	for rows.Next() {
		var entry models.QueueEntry
		if err := rows.Scan(&entry.Position, &entry.ID, &entry.FrontPreview); err != nil {
			return nil, fmt.Errorf("scan queue row: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate queue rows: %w", err)
	}

	return entries, nil
}
