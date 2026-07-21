package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/importexport"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	importDeck   string
	importFormat string
	importFile   string
	importAppend bool
)

var errImportFileRequired = errors.New("import requires --file")

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import cards into a deck from JSON or CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := strings.TrimSpace(importDeck)
		if deckName == "" {
			return models.ErrDeckNameRequired
		}

		format := strings.ToLower(strings.TrimSpace(importFormat))
		if format != "json" && format != "csv" {
			return importexport.ErrInvalidFormat
		}

		filePath := strings.TrimSpace(importFile)
		if filePath == "" {
			return errImportFileRequired
		}

		var reader io.Reader
		if filePath == "-" {
			reader = cmd.InOrStdin()
		} else {
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("open import file: %w", err)
			}
			defer file.Close()
			reader = file
		}

		var cards []importexport.CardInput
		var parseErrors []string

		switch format {
		case "json":
			data, err := importexport.ParseJSON(reader, deckName)
			if err != nil {
				return err
			}
			cards, parseErrors = importexport.CardsFromExport(data)
		case "csv":
			var err error
			cards, parseErrors, err = importexport.ParseCSV(reader)
			if err != nil {
				return err
			}
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			result, err := importCards(ctx, repo, deckName, cards, importAppend)
			if err != nil {
				return handleRepoError(err)
			}
			result.Errors = append(result.Errors, parseErrors...)
			return formatter().PrintImportResult(cmd.OutOrStdout(), result)
		})
	},
}

func importCards(
	ctx context.Context,
	repo *db.Repository,
	deckName string,
	cards []importexport.CardInput,
	append bool,
) (importexport.ImportResult, error) {
	result := importexport.ImportResult{Deck: deckName}

	_, err := repo.GetDeckByName(ctx, deckName)
	if errors.Is(err, db.ErrDeckNotFound) {
		if _, err := repo.CreateDeck(ctx, models.Deck{Name: deckName}); err != nil {
			return result, err
		}
	} else if err != nil {
		return result, err
	} else if !append {
		return result, errDeckAlreadyExists
	}

	for i := len(cards) - 1; i >= 0; i-- {
		card := cards[i]
		if _, err := repo.CreateCard(ctx, deckName, models.Card{
			Front: card.Front,
			Back:  card.Back,
		}); err != nil {
			return result, err
		}
		result.CardsImported++
	}

	return result, nil
}

func init() {
	importCmd.Flags().StringVar(&importDeck, "deck", "", "Target deck name")
	importCmd.Flags().StringVar(&importFormat, "format", "json", "Import format: json or csv")
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "Import file path (- for stdin)")
	importCmd.Flags().BoolVar(&importAppend, "append", false, "Append cards to an existing deck")
	_ = importCmd.MarkFlagRequired("deck")
	_ = importCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(importCmd)
}
