package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePrecedence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	cfgDir := filepath.Join(home, ".config", "cards")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.toml")
	if err := os.WriteFile(cfgPath, []byte(`database = "/from/config.toml"`), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CARDS_DB", "/from/env")
	cfg, err := Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabasePath != "/from/env" {
		t.Fatalf("got %q, want env path", cfg.DatabasePath)
	}
	if cfg.Source != SourceEnv {
		t.Fatalf("got source %q, want env", cfg.Source)
	}

	t.Setenv("CARDS_DB", "")
	cfg, err = Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabasePath != "/from/config.toml" {
		t.Fatalf("got %q, want config file path", cfg.DatabasePath)
	}
	if cfg.Source != SourceConfigFile {
		t.Fatalf("got source %q, want config_file", cfg.Source)
	}

	if err := os.Remove(cfgPath); err != nil {
		t.Fatal(err)
	}
	cfg, err = Resolve()
	if err != nil {
		t.Fatal(err)
	}
	wantDefault := filepath.Join(home, ".local", "share", "cards", "cards.db")
	if cfg.DatabasePath != wantDefault {
		t.Fatalf("got %q, want default %q", cfg.DatabasePath, wantDefault)
	}
	if cfg.Source != SourceDefault {
		t.Fatalf("got source %q, want default", cfg.Source)
	}
}

func TestResolveDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("CARDS_DB", "")

	cfg, err := Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BatchSize != 4 {
		t.Fatalf("got batch_size %d, want 4", cfg.BatchSize)
	}
	if cfg.AgainOffset != 2 {
		t.Fatalf("got again_offset %d, want 2", cfg.AgainOffset)
	}
	if cfg.NudgeThresholdDays != 3 {
		t.Fatalf("got nudge_threshold_days %d, want 3", cfg.NudgeThresholdDays)
	}
}
