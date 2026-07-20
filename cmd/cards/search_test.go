package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func setupSearchTestData(t *testing.T, buf *bytes.Buffer) {
	t.Helper()
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"deck", "create", "spanish", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "What is saudade?",
		"--back", "A deep emotional state of longing.",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "How do you say hello?",
		"--back", "Olá",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"add", "spanish",
		"--front", "What is nostalgia?",
		"--back", "A sentimental longing for the past.",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestSearchJSONPositional(t *testing.T) {
	_, buf := testHarness(t)
	setupSearchTestData(t, buf)

	buf.Reset()
	rootCmd.SetArgs([]string{"search", "saudade", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Cards []models.CardSearchResult `json:"cards"`
		Total int                       `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode search JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Total != 1 || len(resp.Cards) != 1 {
		t.Fatalf("expected 1 result, got total=%d cards=%d", resp.Total, len(resp.Cards))
	}
	if resp.Cards[0].Deck != "portuguese" {
		t.Fatalf("expected portuguese deck, got %q", resp.Cards[0].Deck)
	}
}

func TestSearchJSONTerms(t *testing.T) {
	_, buf := testHarness(t)
	setupSearchTestData(t, buf)

	buf.Reset()
	rootCmd.SetArgs([]string{"search", "--term", "saudade", "--term", "nostalgia", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Cards []models.CardSearchResult `json:"cards"`
		Total int                       `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode search JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Total != 2 {
		t.Fatalf("expected 2 results, got %d", resp.Total)
	}
}

func TestSearchJSONDeckFilter(t *testing.T) {
	_, buf := testHarness(t)
	setupSearchTestData(t, buf)

	buf.Reset()
	rootCmd.SetArgs([]string{"search", "hello", "--deck", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Cards []models.CardSearchResult `json:"cards"`
		Total int                       `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode search JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Total != 1 {
		t.Fatalf("expected 1 result, got %d", resp.Total)
	}
}

func TestSearchEmptyJSON(t *testing.T) {
	_, buf := testHarness(t)
	setupSearchTestData(t, buf)

	buf.Reset()
	rootCmd.SetArgs([]string{"search", "nonexistent", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	if bytes.Contains(buf.Bytes(), []byte(`"cards": null`)) {
		t.Fatalf("expected empty array, got: %s", buf.String())
	}

	var resp struct {
		Cards []models.CardSearchResult `json:"cards"`
		Total int                       `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode search JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Total != 0 || len(resp.Cards) != 0 {
		t.Fatalf("expected empty results, got total=%d cards=%d", resp.Total, len(resp.Cards))
	}
}

func TestSearchMissingTerms(t *testing.T) {
	_, _ = testHarness(t)
	resetCommandFlags(t)
	rootCmd.SetArgs([]string{"search", "--json"})

	err := rootCmd.Execute()
	if !errors.Is(err, errSearchTermsRequired) {
		t.Fatalf("expected errSearchTermsRequired, got %v", err)
	}
}

func TestSearchDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"search", "hello", "--deck", "missing", "--json"})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestSearchTableOutput(t *testing.T) {
	_, buf := testHarness(t)
	setupSearchTestData(t, buf)

	resetCommandFlags(t)
	buf.Reset()
	rootCmd.SetArgs([]string{"search", "saudade"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{"DECK", "ID", "FRONT", "BACK", "portuguese", "saudade"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in table output, got: %s", want, out)
		}
	}
}
