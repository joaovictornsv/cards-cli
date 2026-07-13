package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var deckDeleteYes bool

var deckDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a deck and all its cards",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			if jsonOutput && !deckDeleteYes {
				return fmt.Errorf("delete requires --yes when using --json")
			}

			deck, err := repo.GetDeckByName(ctx, name)
			if err != nil {
				return handleRepoError(err)
			}

			if !deckDeleteYes && isInteractiveTerminal() {
				prompt := fmt.Sprintf(`Delete deck "%s" (%d cards)? [y/N] `, deck.Name, deck.CardCount)
				if _, err := fmt.Fprint(cmd.OutOrStdout(), prompt); err != nil {
					return err
				}
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("read confirmation: %w", err)
				}
				answer := strings.ToLower(strings.TrimSpace(line))
				if answer != "y" && answer != "yes" {
					return fmt.Errorf("delete cancelled")
				}
			}

			deleted, err := repo.DeleteDeckByName(ctx, name)
			if err != nil {
				return handleRepoError(err)
			}
			return formatter().PrintDeck(cmd.OutOrStdout(), deleted)
		})
	},
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func init() {
	deckDeleteCmd.Flags().BoolVarP(&deckDeleteYes, "yes", "y", false, "Confirm deletion without prompting")
	deckCmd.AddCommand(deckDeleteCmd)
}
