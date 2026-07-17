package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	envDatabase   = "CARDS_DB"
	configDirName = "cards"
	configFile    = "config.toml"

	defaultBatchSize   = 4
	defaultAgainOffset = 2
)

type Source string

const (
	SourceEnv        Source = "env"
	SourceConfigFile Source = "config_file"
	SourceDefault    Source = "default"
)

type Config struct {
	DatabasePath string
	ConfigPath   string
	ConfigExists bool
	Source       Source
	BatchSize   int
	AgainOffset int
}

type fileConfig struct {
	Database    string `toml:"database"`
	BatchSize   int    `toml:"batch_size"`
	AgainOffset int    `toml:"again_offset"`
}

func Resolve() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve home directory: %w", err)
	}

	cfgPath := filepath.Join(home, ".config", configDirName, configFile)
	fc := fileConfig{}
	configExists := fileExists(cfgPath)
	if configExists {
		if _, err := toml.DecodeFile(cfgPath, &fc); err != nil {
			return Config{}, fmt.Errorf("read config file %s: %w", cfgPath, err)
		}
	}

	batchSize := defaultBatchSize
	if fc.BatchSize > 0 {
		batchSize = fc.BatchSize
	}
	againOffset := defaultAgainOffset
	if fc.AgainOffset > 0 {
		againOffset = fc.AgainOffset
	}

	if v := os.Getenv(envDatabase); v != "" {
		return Config{
			DatabasePath: v,
			ConfigPath:   cfgPath,
			ConfigExists: configExists,
			Source:       SourceEnv,
			BatchSize:    batchSize,
			AgainOffset:  againOffset,
		}, nil
	}

	if fc.Database != "" {
		return Config{
			DatabasePath: fc.Database,
			ConfigPath:   cfgPath,
			ConfigExists: configExists,
			Source:       SourceConfigFile,
			BatchSize:    batchSize,
			AgainOffset:  againOffset,
		}, nil
	}

	defaultPath := filepath.Join(home, ".local", "share", configDirName, "cards.db")
	return Config{
		DatabasePath: defaultPath,
		ConfigPath:   cfgPath,
		ConfigExists: configExists,
		Source:       SourceDefault,
		BatchSize:    batchSize,
		AgainOffset:  againOffset,
	}, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
