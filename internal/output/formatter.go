package output

import (
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

type Formatter interface {
	PrintConfig(w io.Writer, cfg config.Config) error
	PrintVersion(w io.Writer, info buildinfo.Info) error
	PrintDeck(w io.Writer, deck models.Deck) error
	PrintDecks(w io.Writer, decks []models.Deck) error
	PrintCard(w io.Writer, card models.Card) error
	PrintCards(w io.Writer, deckName string, cards []models.CardSummary) error
	PrintQueue(w io.Writer, deckName string, entries []models.QueueEntry) error
}

func New(json bool) Formatter {
	if json {
		return JSONFormatter{}
	}
	return TableFormatter{}
}
