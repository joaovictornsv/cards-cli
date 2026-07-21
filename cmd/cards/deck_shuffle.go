package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/spf13/cobra"
)

var (
	deckShuffleYes   bool
	deckShuffleSeed  int64
	deckShuffleHasSeed bool
)

var deckShuffleCmd = &cobra.Command{
	Use:   "shuffle [name]",
	Short: "Randomly reshuffle the deck queue order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			if jsonOutput && !deckShuffleYes {
				return errShuffleRequiresYes
			}

			deck, err := repo.GetDeckByName(ctx, name)
			if err != nil {
				return handleRepoError(err)
			}

			if !deckShuffleYes && isInteractiveTerminal() {
				prompt := fmt.Sprintf(`Shuffle queue for deck "%s" (%d cards)? [y/N] `, deck.Name, deck.CardCount)
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
					return fmt.Errorf("shuffle cancelled")
				}
			}

			var rng *rand.Rand
			if deckShuffleHasSeed {
				rng = rand.New(rand.NewSource(deckShuffleSeed))
			} else {
				rng = rand.New(rand.NewSource(time.Now().UnixNano()))
			}

			result, err := repo.ShuffleDeckQueue(ctx, name, rng)
			if err != nil {
				return handleRepoError(err)
			}

			return formatter().PrintShuffleResult(cmd.OutOrStdout(), result)
		})
	},
}

func init() {
	deckShuffleCmd.Flags().BoolVarP(&deckShuffleYes, "yes", "y", false, "Confirm shuffle without prompting")
	deckShuffleCmd.Flags().Int64Var(&deckShuffleSeed, "seed", 0, "Deterministic shuffle seed (for testing)")
	_ = deckShuffleCmd.Flags().MarkHidden("seed")
	deckShuffleCmd.PreRun = func(cmd *cobra.Command, args []string) {
		deckShuffleHasSeed = cmd.Flags().Changed("seed")
	}
	deckCmd.AddCommand(deckShuffleCmd)
}
