package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI version and build metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		return printVersion(cmd.OutOrStdout(), buildinfo.Get())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(w io.Writer, info buildinfo.Info) error {
	if jsonOutput {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}
	_, err := fmt.Fprintf(w, "%s (commit %s, %s)\n", info.Version, info.Commit, info.GoVersion)
	return err
}
