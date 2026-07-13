package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
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

func resetDeckCommandFlags(t *testing.T) {
	t.Helper()
	jsonOutput = false
	deckDeleteYes = false
	rootCmd.SetArgs(nil)
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

func TestDeckCreateDuplicate(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected duplicate error")
	}
	if !strings.Contains(err.Error(), "deck already exists") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeckCreateValidatesBeforeDB(t *testing.T) {
	dbPath, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "   "})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for empty name")
	}
	if !strings.Contains(err.Error(), "deck name is required") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dbPath); statErr == nil {
		t.Fatal("expected database not to be created when validation fails early")
	}
	_ = buf
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

func TestDeckDeleteNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"deck", "delete", "missing", "--json", "--yes"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected not found error")
	}
	if !strings.Contains(err.Error(), "deck not found") {
		t.Fatalf("unexpected error: %v", err)
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
	if err == nil {
		t.Fatal("expected error when --json without --yes")
	}
	if !strings.Contains(err.Error(), "requires --yes") {
		t.Fatalf("unexpected error: %v", err)
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
