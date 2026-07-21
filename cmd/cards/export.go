package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/importexport"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	exportOutput string
)

var exportCmd = &cobra.Command{
	Use:   "export <deck>",
	Short: "Export a deck and its cards to JSON or CSV",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format := strings.ToLower(strings.TrimSpace(exportFormat))
		if format != "json" && format != "csv" {
			return importexport.ErrInvalidFormat
		}

		deckName := args[0]
		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			cards, err := repo.ListCardsInQueueOrder(ctx, deckName)
			if err != nil {
				return handleRepoError(err)
			}

			exports := make([]importexport.CardExport, len(cards))
			for i, card := range cards {
				id := card.ID
				exports[i] = importexport.CardExport{
					ID:    &id,
					Front: card.Front,
					Back:  card.Back,
				}
			}

			outputPath := strings.TrimSpace(exportOutput)
			writePayload := !jsonOutput || outputPath != ""
			if writePayload {
				var out io.Writer = cmd.OutOrStdout()
				if outputPath != "" {
					file, err := os.Create(outputPath)
					if err != nil {
						return fmt.Errorf("create output file: %w", err)
					}
					defer file.Close()
					out = file
				}

				switch format {
				case "json":
					if err := importexport.WriteJSON(out, importexport.DeckExport{
						Deck:  deckName,
						Cards: exports,
					}); err != nil {
						return err
					}
				case "csv":
					if err := importexport.WriteCSV(out, exports); err != nil {
						return err
					}
				}
			}

			if jsonOutput {
				return formatter().PrintExportSummary(cmd.OutOrStdout(), importexport.ExportSummary{
					Deck:      deckName,
					Format:    format,
					CardCount: len(cards),
					Output:    outputPath,
				})
			}
			return nil
		})
	},
}

func init() {
	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Export format: json or csv")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Write export to file instead of stdout")
	rootCmd.AddCommand(exportCmd)
}
