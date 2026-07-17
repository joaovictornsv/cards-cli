package main

import (
	"context"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	editFront           string
	editBack            string
	editReplaceEligible bool
)

var editCmd = &cobra.Command{
	Use:   "edit [deck] [id]",
	Short: "Edit a card's front and/or back",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		cardID, err := parseCardID(args[1])
		if err != nil {
			return err
		}

		var front, back *string
		if cmd.Flags().Changed("front") {
			front = &editFront
		}
		if cmd.Flags().Changed("back") {
			back = &editBack
		}
		var replaceEligible *bool
		if cmd.Flags().Changed("replace-eligible") {
			replaceEligible = &editReplaceEligible
		}
		if front == nil && back == nil && replaceEligible == nil {
			return models.ErrCardEditRequiresField
		}

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			updated, err := repo.UpdateCard(ctx, deckName, cardID, front, back, replaceEligible)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintCard(cmd.OutOrStdout(), updated)
		})
	},
}

func init() {
	editCmd.Flags().StringVar(&editFront, "front", "", "New front text")
	editCmd.Flags().StringVar(&editBack, "back", "", "New back text")
	editCmd.Flags().BoolVar(&editReplaceEligible, "replace-eligible", false, "Set replace_eligible flag (use --replace-eligible=false to clear)")
	rootCmd.AddCommand(editCmd)
}
