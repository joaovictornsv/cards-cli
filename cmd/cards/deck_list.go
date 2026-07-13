package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var deckListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all decks with card counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			decks, err := repo.ListDecks(ctx)
			if err != nil {
				return err
			}
			return formatter().PrintDecks(cmd.OutOrStdout(), decks)
		})
	},
}

func init() {
	deckCmd.AddCommand(deckListCmd)
}
