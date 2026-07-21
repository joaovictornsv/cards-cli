package importexport

import (
	"errors"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

var (
	ErrInvalidFormat    = errors.New("format must be json or csv")
	ErrDeckNameMismatch = errors.New("deck name in file does not match --deck")
	ErrInvalidJSON      = errors.New("invalid JSON")
	ErrInvalidCSVHeader = errors.New("invalid CSV header: expected front,back")
	ErrMissingCSVField  = errors.New("missing required CSV field")
)

// DeckExport is the JSON export shape for a deck and its cards.
type DeckExport struct {
	Deck  string       `json:"deck"`
	Cards []CardExport `json:"cards"`
}

// CardExport is a single card in an export file.
type CardExport struct {
	ID    *int64 `json:"id,omitempty"`
	Front string `json:"front"`
	Back  string `json:"back"`
}

// ImportResult summarizes a completed import operation.
type ImportResult struct {
	Deck          string   `json:"deck"`
	CardsImported int      `json:"cards_imported"`
	Errors        []string `json:"errors,omitempty"`
}

// ExportSummary is the --json response for cards export.
type ExportSummary struct {
	Deck      string `json:"deck"`
	Format    string `json:"format"`
	CardCount int    `json:"card_count"`
	Output    string `json:"output,omitempty"`
}

// CardInput is a validated card ready for database insertion.
type CardInput struct {
	Front string
	Back  string
}

func ValidateCardInput(front, back string) (CardInput, error) {
	card := models.Card{Front: front, Back: back}
	if err := card.ValidateForCreate(); err != nil {
		return CardInput{}, err
	}
	return CardInput{Front: card.Front, Back: card.Back}, nil
}
