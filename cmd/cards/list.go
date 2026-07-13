package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [deck]",
	Short: "List cards in a deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			cards, err := repo.ListCardsByDeck(ctx, deckName)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintCards(cmd.OutOrStdout(), deckName, cards)
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
