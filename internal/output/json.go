package output

import (
	"encoding/json"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/study"
)

type JSONFormatter struct{}

func (JSONFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	payload := map[string]any{
		"database_path": cfg.DatabasePath,
		"config_path":   cfg.ConfigPath,
		"config_exists": cfg.ConfigExists,
		"source":        cfg.Source,
		"batch_size":   cfg.BatchSize,
		"again_offset": cfg.AgainOffset,
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

type searchResponse struct {
	Cards []models.CardSearchResult `json:"cards"`
	Total int                       `json:"total"`
}

func (JSONFormatter) PrintSearchResults(w io.Writer, results []models.CardSearchResult) error {
	if results == nil {
		results = []models.CardSearchResult{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(searchResponse{
		Cards: results,
		Total: len(results),
	})
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

type queueResponse struct {
	Deck  string              `json:"deck"`
	Queue []models.QueueEntry `json:"queue"`
}

func (JSONFormatter) PrintQueue(w io.Writer, deckName string, entries []models.QueueEntry) error {
	if entries == nil {
		entries = []models.QueueEntry{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(queueResponse{
		Deck:  deckName,
		Queue: entries,
	})
}

func (JSONFormatter) PrintStudyLog(w io.Writer, result study.Result) error {
	if result.Reviews == nil {
		result.Reviews = []study.Review{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
