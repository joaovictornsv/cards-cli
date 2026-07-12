package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
)

func TestJSONFormatter(t *testing.T) {
	cfg := config.Config{
		DatabasePath: "/home/user/.local/share/cards/cards.db",
		ConfigPath:   "/home/user/.config/cards/config.toml",
		ConfigExists: false,
		Source:       config.SourceDefault,
		BatchSize:    4,
		AgainOffset:  2,
		HardOffset:   5,
	}

	var buf bytes.Buffer
	formatter := JSONFormatter{}
	if err := formatter.PrintConfig(&buf, cfg); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		`"database_path": "/home/user/.local/share/cards/cards.db"`,
		`"config_path": "/home/user/.config/cards/config.toml"`,
		`"config_exists": false`,
		`"source": "default"`,
		`"batch_size": 4`,
		`"again_offset": 2`,
		`"hard_offset": 5`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in config json, got: %s", want, out)
		}
	}

	buf.Reset()
	info := buildinfo.Info{
		Version:   "0.0.0-dev",
		Commit:    "unknown",
		GoVersion: "go1.25.0",
	}
	if err := formatter.PrintVersion(&buf, info); err != nil {
		t.Fatal(err)
	}
	out = buf.String()
	for _, want := range []string{
		`"version": "0.0.0-dev"`,
		`"commit": "unknown"`,
		`"go_version": "go1.25.0"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in version json, got: %s", want, out)
		}
	}
}

func TestTableFormatter(t *testing.T) {
	cfg := config.Config{
		DatabasePath: "/home/user/.local/share/cards/cards.db",
		ConfigPath:   "/home/user/.config/cards/config.toml",
		ConfigExists: false,
		Source:       config.SourceDefault,
		BatchSize:    4,
		AgainOffset:  2,
		HardOffset:   5,
	}

	var buf bytes.Buffer
	table := TableFormatter{}
	if err := table.PrintConfig(&buf, cfg); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"database_path: /home/user/.local/share/cards/cards.db",
		"config_path: /home/user/.config/cards/config.toml",
		"config_exists: false",
		"source: default",
		"batch_size: 4",
		"again_offset: 2",
		"hard_offset: 5",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in config table, got: %s", want, out)
		}
	}

	buf.Reset()
	info := buildinfo.Info{
		Version:   "0.0.0-dev",
		Commit:    "unknown",
		GoVersion: "go1.25.0",
	}
	if err := table.PrintVersion(&buf, info); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "0.0.0-dev (commit unknown, go1.25.0)") {
		t.Fatalf("unexpected version table: %s", buf.String())
	}
}
