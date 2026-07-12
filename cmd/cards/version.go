package main

import (
	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version and build metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		return formatter().PrintVersion(cmd.OutOrStdout(), buildinfo.Get())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
