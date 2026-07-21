package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/queue"
	"github.com/joaovictornsv/cards-cli/internal/study"
	"github.com/spf13/cobra"
)

var (
	studyLimit      int
	studyInputFactory = func(in io.Reader) study.Input {
		return study.NewTerminalInput(in)
	}
)

var studyCmd = &cobra.Command{
	Use:   "study [deck]",
	Short: "Run an interactive study session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deckName := args[0]
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}

		batchSize := cfg.BatchSize
		if cmd.Flags().Changed("limit") {
			if studyLimit < 1 {
				return fmt.Errorf("--limit must be at least 1")
			}
			batchSize = studyLimit
		}

		return runStudySession(cmd.Context(), deckName, cfg, batchSize, cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func runStudySession(ctx context.Context, deckName string, cfg config.Config, batchSize int, in io.Reader, out io.Writer) error {
	repo, cleanup, err := openRepo()
	if err != nil {
		return err
	}
	defer cleanup()

	return runStudyWithRepo(ctx, repo, deckName, cfg, batchSize, in, out)
}

func runStudyWithRepo(ctx context.Context, repo *db.Repository, deckName string, cfg config.Config, batchSize int, in io.Reader, out io.Writer) error {
	sess := &study.Session{
		DeckName: deckName,
		Out:      out,
		Store:    study.NewDBStore(repo),
		Input:    studyInputFactory(in),
		Opts: study.Options{
			BatchSize: batchSize,
			QueueOpts: queue.Options{
				AgainOffset: cfg.AgainOffset,
			},
		},
	}

	result, err := sess.Run(ctx)
	if err != nil {
		if errors.Is(err, study.ErrDeckNotFound) {
			return errDeckNotFound
		}
		if errors.Is(err, study.ErrEmptyDeck) {
			return fmt.Errorf(
				`deck %q has no cards — add cards with: cards add %s --front "..." --back "..."`,
				deckName, deckName,
			)
		}
		return err
	}

	if shouldRecordSession(result) {
		if err := repo.RecordDeckSessionByName(ctx, deckName, models.NowTimestamp()); err != nil {
			return err
		}
	}

	if jsonOutput {
		fmt.Fprintln(out)
		if err := formatter().PrintStudyLog(out, result); err != nil {
			return err
		}
	}

	return nil
}

func shouldRecordSession(result study.Result) bool {
	if result.Status == "complete" {
		return true
	}
	if result.Status == "quit" && len(result.Reviews) > 0 {
		return true
	}
	return false
}

func init() {
	studyCmd.Flags().IntVar(&studyLimit, "limit", 0, "Batch size override (default from config)")
	rootCmd.AddCommand(studyCmd)
}
