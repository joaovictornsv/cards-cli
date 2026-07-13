package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	addFront string
	addBack  string
)

var addCmd = &cobra.Command{
	Use:   "add [deck]",
	Short: "Add a card to a deck",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		card := models.Card{Front: addFront, Back: addBack}
		if err := card.ValidateForCreate(); err != nil {
			return err
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			created, err := repo.CreateCard(ctx, args[0], card)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintCard(cmd.OutOrStdout(), created)
		})
	},
}

func init() {
	addCmd.Flags().StringVar(&addFront, "front", "", "Card front text")
	addCmd.Flags().StringVar(&addBack, "back", "", "Card back text")
	_ = addCmd.MarkFlagRequired("front")
	_ = addCmd.MarkFlagRequired("back")
	rootCmd.AddCommand(addCmd)
}
