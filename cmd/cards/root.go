package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/output"
	"github.com/spf13/cobra"
)

var jsonOutput bool

var rootCmd = &cobra.Command{
	Use:           "cards",
	Short:         "CLI flashcard app for terminal study sessions",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Machine-readable JSON output")
	rootCmd.Version = buildinfo.Version
}

func openRepo() (*db.Repository, func(), error) {
	cfg, err := config.Resolve()
	if err != nil {
		return nil, nil, err
	}

	database, err := db.Open(cfg.DatabasePath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = database.Close()
	}
	return db.NewRepository(database), cleanup, nil
}

func formatter() output.Formatter {
	return output.New(jsonOutput)
}

func handleRepoError(err error) error {
	if errors.Is(err, db.ErrNotFound) {
		return fmt.Errorf("deck not found")
	}
	if errors.Is(err, db.ErrDuplicateName) {
		return fmt.Errorf("deck already exists")
	}
	return err
}

func runWithRepo(ctx context.Context, fn func(context.Context, *db.Repository) error) error {
	repo, cleanup, err := openRepo()
	if err != nil {
		return err
	}
	defer cleanup()
	return fn(ctx, repo)
}
