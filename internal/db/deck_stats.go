package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type deckStatsRow struct {
	Deck          string
	SessionsCount int
	LastSessionAt sql.NullString
}

func (r *Repository) GetDeckStatsByName(ctx context.Context, name string) (deckStatsRow, error) {
	row := r.db.sql.QueryRowContext(ctx, `
		SELECT name, sessions_count, last_session_at
		FROM decks
		WHERE name = ?`, name)

	var stats deckStatsRow
	if err := row.Scan(&stats.Deck, &stats.SessionsCount, &stats.LastSessionAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deckStatsRow{}, ErrDeckNotFound
		}
		return deckStatsRow{}, fmt.Errorf("get deck stats: %w", err)
	}
	return stats, nil
}

func (r *Repository) RecordDeckSession(ctx context.Context, deckID int64, at string) error {
	res, err := r.db.sql.ExecContext(ctx, `
		UPDATE decks
		SET sessions_count = sessions_count + 1,
		    last_session_at = ?
		WHERE id = ?`, at, deckID)
	if err != nil {
		return fmt.Errorf("record deck session: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrDeckNotFound
	}
	return nil
}

func (r *Repository) RecordDeckSessionByName(ctx context.Context, deckName string, at string) error {
	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return err
	}
	return r.RecordDeckSession(ctx, deck.ID, at)
}

func DeckStatsRowToModel(row deckStatsRow) (string, int, *string) {
	var lastSessionAt *string
	if row.LastSessionAt.Valid && row.LastSessionAt.String != "" {
		v := row.LastSessionAt.String
		lastSessionAt = &v
	}
	return row.Deck, row.SessionsCount, lastSessionAt
}
