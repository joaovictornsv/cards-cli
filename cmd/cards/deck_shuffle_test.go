package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func createDeckWithNCards(t *testing.T, deckName string, count int) {
	t.Helper()
	rootCmd.SetArgs([]string{"deck", "create", deckName, "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < count; i++ {
		rootCmd.SetArgs([]string{
			"add", deckName,
			"--front", "front",
			"--back", "back",
			"--json",
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeckShuffleJSON(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithNCards(t, "portuguese", 3)

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "shuffle", "portuguese", "--json", "--yes", "--seed", "42"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result models.ShuffleResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("decode shuffle JSON: %v\noutput: %s", err, buf.String())
	}
	if result.Deck != "portuguese" {
		t.Fatalf("expected deck portuguese, got %q", result.Deck)
	}
	if result.CardCount != 3 {
		t.Fatalf("expected card_count 3, got %d", result.CardCount)
	}
	if result.Status != "shuffled" {
		t.Fatalf("expected status shuffled, got %q", result.Status)
	}
}

func TestDeckShuffleRequiresYesWithJSON(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithNCards(t, "portuguese", 3)

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "shuffle", "portuguese", "--json"})
	err := rootCmd.Execute()
	if !errors.Is(err, errShuffleRequiresYes) {
		t.Fatalf("expected errShuffleRequiresYes, got %v", err)
	}
}

func TestDeckShuffleNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"deck", "shuffle", "missing", "--json", "--yes"})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestDeckShuffleNoopSingleCard(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithNCards(t, "portuguese", 1)

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "shuffle", "portuguese", "--json", "--yes"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result models.ShuffleResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("decode shuffle JSON: %v\noutput: %s", err, buf.String())
	}
	if result.Status != "noop" || result.CardCount != 1 {
		t.Fatalf("expected noop with 1 card, got %+v", result)
	}
}

func TestDeckShuffleNoopEmptyDeck(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "empty", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"deck", "shuffle", "empty", "--json", "--yes"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result models.ShuffleResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("decode shuffle JSON: %v\noutput: %s", err, buf.String())
	}
	if result.Status != "noop" || result.CardCount != 0 {
		t.Fatalf("expected noop with 0 cards, got %+v", result)
	}
}
