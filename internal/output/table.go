package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

type TableFormatter struct{}

func (TableFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	return PrintConfigHuman(w, cfg)
}

func (TableFormatter) PrintVersion(w io.Writer, info buildinfo.Info) error {
	return PrintVersionHuman(w, info)
}

func PrintConfigHuman(w io.Writer, cfg config.Config) error {
	_, err := fmt.Fprintf(w,
		"database_path: %s\nconfig_path: %s\nconfig_exists: %t\nsource: %s\nbatch_size: %d\nagain_offset: %d\nhard_offset: %d\n",
		cfg.DatabasePath, cfg.ConfigPath, cfg.ConfigExists, cfg.Source,
		cfg.BatchSize, cfg.AgainOffset, cfg.HardOffset)
	return err
}

func PrintVersionHuman(w io.Writer, info buildinfo.Info) error {
	_, err := fmt.Fprintf(w, "%s (commit %s, %s)\n", info.Version, info.Commit, info.GoVersion)
	return err
}

func (TableFormatter) PrintDeck(w io.Writer, deck models.Deck) error {
	return printDecksTable(w, []models.Deck{deck})
}

func (TableFormatter) PrintDecks(w io.Writer, decks []models.Deck) error {
	return printDecksTable(w, decks)
}

func printDecksTable(w io.Writer, decks []models.Deck) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "ID\tNAME\tCARDS\tCREATED"); err != nil {
		return err
	}
	for _, deck := range decks {
		if _, err := fmt.Fprintf(tw, "%d\t%s\t%d\t%s\n",
			deck.ID, deck.Name, deck.CardCount, deck.CreatedAt); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (TableFormatter) PrintCard(w io.Writer, card models.Card) error {
	_, err := fmt.Fprintf(w,
		"id: %d\nfront: %s\nback: %s\ncreated_at: %s\nupdated_at: %s\n",
		card.ID, card.Front, card.Back, card.CreatedAt, card.UpdatedAt,
	)
	return err
}

func (TableFormatter) PrintCards(w io.Writer, deckName string, cards []models.CardSummary) error {
	if _, err := fmt.Fprintf(w, "deck: %s\n", deckName); err != nil {
		return err
	}
	return printCardsTable(w, cards)
}

func printCardsTable(w io.Writer, cards []models.CardSummary) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "ID\tFRONT\tCREATED\tUPDATED"); err != nil {
		return err
	}
	for _, card := range cards {
		if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
			card.ID, card.Front, card.CreatedAt, card.UpdatedAt); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (TableFormatter) PrintQueue(w io.Writer, deckName string, entries []models.QueueEntry) error {
	if _, err := fmt.Fprintf(w, "deck: %s\n", deckName); err != nil {
		return err
	}
	return printQueueTable(w, entries)
}

func printQueueTable(w io.Writer, entries []models.QueueEntry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "POSITION\tID\tFRONT"); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := fmt.Fprintf(tw, "%d\t%d\t%s\n",
			entry.Position, entry.ID, entry.FrontPreview); err != nil {
			return err
		}
	}
	return tw.Flush()
}
