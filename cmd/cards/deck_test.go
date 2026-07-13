package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func testHarness(t *testing.T) (string, *bytes.Buffer) {
	t.Helper()
	resetDeckCommandFlags(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	dbPath := filepath.Join(home, "cards.db")
	t.Setenv("CARDS_DB", dbPath)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	return dbPath, &buf
}

func resetCommandFlags(t *testing.T) {
	t.Helper()
	jsonOutput = false
	deckDeleteYes = false
	addFront = ""
	addBack = ""
	editFront = ""
	editBack = ""
	resetCmdFlags(editCmd)
	rootCmd.SetArgs(nil)
}

func resetCmdFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
}

func resetDeckCommandFlags(t *testing.T) {
	resetCommandFlags(t)
}

func TestDeckCreateJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var deck models.Deck
	if err := json.Unmarshal(buf.Bytes(), &deck); err != nil {
		t.Fatalf("decode deck JSON: %v\noutput: %s", err, buf.String())
	}
	if deck.Name != "portuguese" {
		t.Fatalf("expected name portuguese, got %q", deck.Name)
	}
	if deck.ID <= 0 {
		t.Fatalf("expected positive id, got %d", deck.ID)
	}
}

func TestDeckCreateTrimsWhitespace(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "  portuguese  ", "--json"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var deck models.Deck
	if err := json.Unmarshal(buf.Bytes(), &deck); err != nil {
		t.Fatalf("decode deck JSON: %v\noutput: %s", err, buf.String())
	}
	if deck.Name != "portuguese" {
		t.Fatalf("expected trimmed name portuguese, got %q", deck.Name)
	}
}

func TestDeckCreateDuplicate(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeckAlreadyExists) {
		t.Fatalf("expected errDeckAlreadyExists, got %v", err)
	}
}

func TestDeckCreateValidatesBeforeDB(t *testing.T) {
	dbPath, _ := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "   "})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
	if !errors.Is(err, models.ErrDeckNameRequired) {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dbPath); statErr == nil {
		t.Fatal("expected database not to be created when validation fails early")
	}
}

func TestDeckListJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Decks []models.Deck `json:"decks"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode list JSON: %v\noutput: %s", err, buf.String())
	}
	if len(resp.Decks) != 1 {
		t.Fatalf("expected 1 deck, got %d", len(resp.Decks))
	}
	if resp.Decks[0].Name != "portuguese" {
		t.Fatalf("expected name portuguese, got %q", resp.Decks[0].Name)
	}
	if resp.Decks[0].CardCount != 0 {
		t.Fatalf("expected card_count 0, got %d", resp.Decks[0].CardCount)
	}
}

func TestDeckListEmptyJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "list", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	raw := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`"decks": null`)) {
		t.Fatalf("expected empty array, got: %s", raw)
	}

	var resp struct {
		Decks []models.Deck `json:"decks"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode list JSON: %v\noutput: %s", err, raw)
	}
	if len(resp.Decks) != 0 {
		t.Fatalf("expected 0 decks, got %d", len(resp.Decks))
	}
}

func TestDeckDeleteNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"deck", "delete", "missing", "--json", "--yes"})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestDeckDeleteRequiresYesWithJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "delete", "portuguese", "--json"})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeleteRequiresYes) {
		t.Fatalf("expected errDeleteRequiresYes, got %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "delete", "portuguese", "--json", "--yes"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var deck models.Deck
	if err := json.Unmarshal(buf.Bytes(), &deck); err != nil {
		t.Fatalf("decode deleted deck JSON: %v\noutput: %s", err, buf.String())
	}
	if deck.Name != "portuguese" {
		t.Fatalf("expected name portuguese, got %q", deck.Name)
	}
}
