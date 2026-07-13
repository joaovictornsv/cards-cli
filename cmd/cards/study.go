package main

import (
	"context"
	"errors"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/queue"
	"github.com/joaovictornsv/cards-cli/internal/study"
	"github.com/spf13/cobra"
)

var studyInputFactory = func(in io.Reader) study.Input {
	return study.NewTerminalInput(in)
}

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
		return runStudySession(cmd.Context(), deckName, cfg, cmd.InOrStdin(), cmd.OutOrStdout())
	},
}

func runStudySession(ctx context.Context, deckName string, cfg config.Config, in io.Reader, out io.Writer) error {
	repo, cleanup, err := openRepo()
	if err != nil {
		return err
	}
	defer cleanup()

	return runStudyWithRepo(ctx, repo, deckName, cfg, in, out)
}

func runStudyWithRepo(ctx context.Context, repo *db.Repository, deckName string, cfg config.Config, in io.Reader, out io.Writer) error {
	sess := &study.Session{
		DeckName: deckName,
		Out:      out,
		Store:    study.NewDBStore(repo),
		Input:    studyInputFactory(in),
		Opts: study.Options{
			BatchSize: cfg.BatchSize,
			QueueOpts: queue.Options{
				AgainOffset: cfg.AgainOffset,
				HardOffset:  cfg.HardOffset,
			},
		},
	}

	if err := sess.Run(ctx); err != nil {
		if errors.Is(err, study.ErrDeckNotFound) {
			return errDeckNotFound
		}
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(studyCmd)
}
