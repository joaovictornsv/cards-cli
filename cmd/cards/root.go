package main

import (
	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
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
