package db

import (
	"context"
	"fmt"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

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
