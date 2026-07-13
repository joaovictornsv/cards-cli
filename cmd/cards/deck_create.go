package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/spf13/cobra"
)

var deckCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deck := models.Deck{Name: args[0]}
		if err := deck.ValidateForCreate(); err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			created, err := repo.CreateDeck(ctx, deck)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintDeck(cmd.OutOrStdout(), created)
		})
	},
}

func init() {
	deckCmd.AddCommand(deckCreateCmd)
}
