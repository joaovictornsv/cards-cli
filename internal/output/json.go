package output

import (
	"encoding/json"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

type JSONFormatter struct{}

func (JSONFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	payload := map[string]any{
		"database_path": cfg.DatabasePath,
		"config_path":   cfg.ConfigPath,
		"config_exists": cfg.ConfigExists,
		"source":        cfg.Source,
		"batch_size":    cfg.BatchSize,
		"again_offset":  cfg.AgainOffset,
		"hard_offset":   cfg.HardOffset,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}

func (JSONFormatter) PrintVersion(w io.Writer, info buildinfo.Info) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(info)
}

func (JSONFormatter) PrintDeck(w io.Writer, deck models.Deck) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(deck)
}

type decksResponse struct {
	Decks []models.Deck `json:"decks"`
}

func (JSONFormatter) PrintDecks(w io.Writer, decks []models.Deck) error {
	if decks == nil {
		decks = []models.Deck{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(decksResponse{Decks: decks})
}

func (JSONFormatter) PrintCard(w io.Writer, card models.Card) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(card)
}

type cardsResponse struct {
	Deck  string               `json:"deck"`
	Cards []models.CardSummary `json:"cards"`
	Total int                  `json:"total"`
}

func (JSONFormatter) PrintCards(w io.Writer, deckName string, cards []models.CardSummary) error {
	if cards == nil {
		cards = []models.CardSummary{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(cardsResponse{
		Deck:  deckName,
		Cards: cards,
		Total: len(cards),
	})
}
