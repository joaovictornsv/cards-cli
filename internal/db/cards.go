package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

var ErrCardNotFound = errors.New("card not found")

const cardSelectBase = `
	SELECT id, deck_id, front, back, created_at, updated_at, replace_eligible
	FROM cards`

func (r *Repository) CreateCard(ctx context.Context, deckName string, card models.Card) (models.Card, error) {
	if err := card.ValidateForCreate(); err != nil {
		return models.Card{}, err
	}

	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return models.Card{}, err
	}

	now := models.NowTimestamp()
	if card.CreatedAt == "" {
		card.CreatedAt = now
	}
	if card.UpdatedAt == "" {
		card.UpdatedAt = now
	}
	card.DeckID = deck.ID

	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return models.Card{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO cards (deck_id, front, back, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`,
		card.DeckID, card.Front, card.Back, card.CreatedAt, card.UpdatedAt,
	)
	if err != nil {
		return models.Card{}, fmt.Errorf("insert card: %w", err)
	}

	cardID, err := res.LastInsertId()
	if err != nil {
		return models.Card{}, fmt.Errorf("last insert id: %w", err)
	}

	if err := shiftQueueForFrontInsert(ctx, tx, deck.ID); err != nil {
		return models.Card{}, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO queue (deck_id, position, card_id)
		VALUES (?, 0, ?)`,
		deck.ID, cardID,
	); err != nil {
		return models.Card{}, fmt.Errorf("insert queue entry: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return models.Card{}, fmt.Errorf("commit transaction: %w", err)
	}

	return r.GetCardByID(ctx, cardID)
}

func (r *Repository) ListCardsByDeck(ctx context.Context, deckName string, replaceEligibleOnly bool) ([]models.CardSummary, error) {
	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, front, created_at, updated_at, replace_eligible
		FROM cards
		WHERE deck_id = ?`
	if replaceEligibleOnly {
		query += ` AND replace_eligible = 1`
	}
	query += ` ORDER BY id`

	rows, err := r.db.sql.QueryContext(ctx, query, deck.ID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()

	cards := make([]models.CardSummary, 0)
	for rows.Next() {
		var card models.CardSummary
		if err := rows.Scan(&card.ID, &card.Front, &card.CreatedAt, &card.UpdatedAt, &card.ReplaceEligible); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *Repository) GetCardByID(ctx context.Context, id int64) (models.Card, error) {
	row := r.db.sql.QueryRowContext(ctx, cardSelectBase+` WHERE id = ?`, id)

	card, err := scanCard(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Card{}, ErrCardNotFound
	}
	if err != nil {
		return models.Card{}, err
	}
	return card, nil
}

func (r *Repository) GetCardByDeckAndID(ctx context.Context, deckName string, cardID int64) (models.Card, error) {
	deck, err := r.GetDeckByName(ctx, deckName)
	if err != nil {
		return models.Card{}, err
	}

	card, err := r.GetCardByID(ctx, cardID)
	if err != nil {
		return models.Card{}, err
	}
	if card.DeckID != deck.ID {
		return models.Card{}, ErrCardNotFound
	}
	return card, nil
}

func (r *Repository) UpdateCard(ctx context.Context, deckName string, cardID int64, front, back *string, replaceEligible *bool) (models.Card, error) {
	card, err := r.GetCardByDeckAndID(ctx, deckName, cardID)
	if err != nil {
		return models.Card{}, err
	}

	if err := card.ValidateForUpdate(front, back, replaceEligible); err != nil {
		return models.Card{}, err
	}

	if front != nil {
		card.Front = *front
	}
	if back != nil {
		card.Back = *back
	}
	if replaceEligible != nil {
		card.ReplaceEligible = *replaceEligible
	}
	card.UpdatedAt = models.NowTimestamp()

	_, err = r.db.sql.ExecContext(ctx, `
		UPDATE cards SET front = ?, back = ?, updated_at = ?, replace_eligible = ?
		WHERE id = ?`,
		card.Front, card.Back, card.UpdatedAt, card.ReplaceEligible, card.ID,
	)
	if err != nil {
		return models.Card{}, fmt.Errorf("update card: %w", err)
	}

	return card, nil
}

func (r *Repository) DeleteCard(ctx context.Context, deckName string, cardID int64) (models.Card, error) {
	card, err := r.GetCardByDeckAndID(ctx, deckName, cardID)
	if err != nil {
		return models.Card{}, err
	}

	tx, err := r.db.sql.BeginTx(ctx, nil)
	if err != nil {
		return models.Card{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var deletedPosition int
	err = tx.QueryRowContext(ctx, `
		SELECT position FROM queue WHERE deck_id = ? AND card_id = ?`,
		card.DeckID, card.ID,
	).Scan(&deletedPosition)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Card{}, ErrCardNotFound
	}
	if err != nil {
		return models.Card{}, fmt.Errorf("get queue position: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM cards WHERE id = ?`, card.ID); err != nil {
		return models.Card{}, fmt.Errorf("delete card: %w", err)
	}

	if err := compactQueueAfterDelete(ctx, tx, card.DeckID, deletedPosition); err != nil {
		return models.Card{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Card{}, fmt.Errorf("commit transaction: %w", err)
	}

	return card, nil
}

// compactQueueAfterDelete renumbers queue positions after a card is removed.
// Uses temporary negative positions to avoid primary-key collisions on
// (deck_id, position).
func compactQueueAfterDelete(ctx context.Context, tx *sql.Tx, deckID int64, deletedPosition int) error {
	if _, err := tx.ExecContext(ctx, `
		UPDATE queue SET position = -(position - ?)
		WHERE deck_id = ? AND position > ?`,
		deletedPosition, deckID, deletedPosition,
	); err != nil {
		return fmt.Errorf("shift queue positions after delete: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE queue SET position = ? - position - 1
		WHERE deck_id = ? AND position < 0`,
		deletedPosition, deckID,
	); err != nil {
		return fmt.Errorf("normalize queue positions after delete: %w", err)
	}
	return nil
}

// shiftQueueForFrontInsert moves existing queue entries back by one position so
// position 0 is free. Uses temporary negative positions to avoid primary-key
// collisions on (deck_id, position).
func shiftQueueForFrontInsert(ctx context.Context, tx *sql.Tx, deckID int64) error {
	if _, err := tx.ExecContext(ctx, `
		UPDATE queue SET position = -(position + 1)
		WHERE deck_id = ?`,
		deckID,
	); err != nil {
		return fmt.Errorf("shift queue positions: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE queue SET position = -position
		WHERE deck_id = ?`,
		deckID,
	); err != nil {
		return fmt.Errorf("normalize queue positions: %w", err)
	}
	return nil
}

func (r *Repository) SetReplaceEligible(ctx context.Context, deckName string, cardID int64, eligible bool) error {
	card, err := r.GetCardByDeckAndID(ctx, deckName, cardID)
	if err != nil {
		return err
	}

	_, err = r.db.sql.ExecContext(ctx, `
		UPDATE cards SET replace_eligible = ?, updated_at = ?
		WHERE id = ?`,
		eligible, models.NowTimestamp(), card.ID,
	)
	if err != nil {
		return fmt.Errorf("set replace_eligible: %w", err)
	}
	return nil
}

func scanCard(row rowScanner) (models.Card, error) {
	var card models.Card
	if err := row.Scan(
		&card.ID, &card.DeckID, &card.Front, &card.Back,
		&card.CreatedAt, &card.UpdatedAt, &card.ReplaceEligible,
	); err != nil {
		return models.Card{}, err
	}
	return card, nil
}
