package main

import (
	"context"
	"time"

	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [deck]",
	Short: "Show deck study stats and session nudge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		return runWithRepo(cmd.Context(), func(ctx context.Context, repo *db.Repository) error {
			return printDeckStats(cmd, repo, deckName, cfg)
		})
	},
}

func printDeckStats(cmd *cobra.Command, repo *db.Repository, deckName string, cfg config.Config) error {
	row, err := repo.GetDeckStatsByName(cmd.Context(), deckName)
	if err != nil {
		return handleRepoError(err)
	}

	deck, sessionsCount, lastSessionAt := db.DeckStatsRowToModel(row)
	view := stats.BuildDeckStats(deck, sessionsCount, lastSessionAt, cfg.NudgeThresholdDays, time.Now().UTC())
	return formatter().PrintDeckStats(cmd.OutOrStdout(), view)
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
