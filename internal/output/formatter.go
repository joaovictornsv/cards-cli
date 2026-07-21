package output

import (
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/importexport"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/study"
)

type Formatter interface {
	PrintConfig(w io.Writer, cfg config.Config) error
	PrintVersion(w io.Writer, info buildinfo.Info) error
	PrintDeck(w io.Writer, deck models.Deck) error
	PrintDecks(w io.Writer, decks []models.Deck) error
	PrintCard(w io.Writer, card models.Card) error
	PrintCards(w io.Writer, deckName string, cards []models.CardSummary) error
	PrintSearchResults(w io.Writer, results []models.CardSearchResult) error
	PrintQueue(w io.Writer, deckName string, entries []models.QueueEntry) error
	PrintStudyLog(w io.Writer, result study.Result) error
	PrintExportSummary(w io.Writer, summary importexport.ExportSummary) error
	PrintImportResult(w io.Writer, result importexport.ImportResult) error
	PrintDeckStats(w io.Writer, stats models.DeckStats) error
}

func New(json bool) Formatter {
	if json {
		return JSONFormatter{}
	}
	return TableFormatter{}
}
