package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var cardDeleteCmd = &cobra.Command{
	Use:   "delete [deck] [id]",
	Short: "Remove a card from a deck and queue",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		cardID, err := parseCardID(args[1])
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			deleted, err := repo.DeleteCard(ctx, deckName, cardID)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintCard(cmd.OutOrStdout(), deleted)
		})
	},
}

func init() {
	rootCmd.AddCommand(cardDeleteCmd)
}
