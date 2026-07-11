package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show effective CLI configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Resolve()
		if err != nil {
			return err
		}
		return printConfig(cmd.OutOrStdout(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func printConfig(w io.Writer, cfg config.Config) error {
	if jsonOutput {
		payload := map[string]any{
			"database_path": cfg.DatabasePath,
			"config_path":   cfg.ConfigPath,
			"config_exists": cfg.ConfigExists,
			"source":        cfg.Source,
			"batch_size":    cfg.BatchSize,
			"again_offset":  cfg.AgainOffset,
			"hard_offset":   cfg.HardOffset,
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}
	_, err := fmt.Fprintf(w,
		"database_path: %s\nconfig_path: %s\nconfig_exists: %t\nsource: %s\nbatch_size: %d\nagain_offset: %d\nhard_offset: %d\n",
		cfg.DatabasePath, cfg.ConfigPath, cfg.ConfigExists, cfg.Source,
		cfg.BatchSize, cfg.AgainOffset, cfg.HardOffset)
	return err
}
