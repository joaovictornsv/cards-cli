package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var queueCmd = &cobra.Command{
	Use:   "queue [deck]",
	Short: "Show current queue order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			entries, err := repo.ListQueueByDeck(ctx, deckName)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintQueue(cmd.OutOrStdout(), deckName, entries)
		})
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)
}
