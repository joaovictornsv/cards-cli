package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [deck] [id]",
	Short: "Show one card",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		cardID, err := parseCardID(args[1])
		if err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			card, err := repo.GetCardByDeckAndID(ctx, deckName, cardID)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintCard(cmd.OutOrStdout(), card)
		})
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
