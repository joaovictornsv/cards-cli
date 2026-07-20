package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func (r *Repository) SearchCards(ctx context.Context, terms []string, deckName string) ([]models.CardSearchResult, error) {
	if len(terms) == 0 {
		return nil, fmt.Errorf("search requires at least one term")
	}

	var deckID *int64
	if deckName != "" {
		deck, err := r.GetDeckByName(ctx, deckName)
		if err != nil {
			return nil, err
		}
		deckID = &deck.ID
	}

	conditions := make([]string, 0, len(terms))
	args := make([]any, 0, len(terms)*3+1)
	for _, term := range terms {
		pattern := "%" + escapeLike(strings.ToLower(term)) + "%"
		conditions = append(conditions, `(
			LOWER(c.front) LIKE ? ESCAPE '\'
			OR LOWER(c.back) LIKE ? ESCAPE '\'
			OR LOWER(d.name) LIKE ? ESCAPE '\'
		)`)
		args = append(args, pattern, pattern, pattern)
	}

	query := `
		SELECT c.id, d.name, c.front, c.back
		FROM cards c
		JOIN decks d ON d.id = c.deck_id
		WHERE (` + strings.Join(conditions, " OR ") + `)`

	if deckID != nil {
		query += ` AND c.deck_id = ?`
		args = append(args, *deckID)
	}
	query += ` ORDER BY d.name, c.id`

	rows, err := r.db.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("search cards: %w", err)
	}
	defer rows.Close()

	results := make([]models.CardSearchResult, 0)
	for rows.Next() {
		var result models.CardSearchResult
		if err := rows.Scan(&result.ID, &result.Deck, &result.Front, &result.Back); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
