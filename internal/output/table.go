package output

import (
	"fmt"
	"io"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
)

type TableFormatter struct{}

func (TableFormatter) PrintConfig(w io.Writer, cfg config.Config) error {
	return PrintConfigHuman(w, cfg)
}

func (TableFormatter) PrintVersion(w io.Writer, info buildinfo.Info) error {
	return PrintVersionHuman(w, info)
}

func PrintConfigHuman(w io.Writer, cfg config.Config) error {
	_, err := fmt.Fprintf(w,
		"database_path: %s\nconfig_path: %s\nconfig_exists: %t\nsource: %s\nbatch_size: %d\nagain_offset: %d\nhard_offset: %d\n",
		cfg.DatabasePath, cfg.ConfigPath, cfg.ConfigExists, cfg.Source,
		cfg.BatchSize, cfg.AgainOffset, cfg.HardOffset)
	return err
}

func PrintVersionHuman(w io.Writer, info buildinfo.Info) error {
	_, err := fmt.Fprintf(w, "%s (commit %s, %s)\n", info.Version, info.Commit, info.GoVersion)
	return err
}
