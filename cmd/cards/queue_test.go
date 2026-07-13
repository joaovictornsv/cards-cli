package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestQueueJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "first card",
		"--back", "first back",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var first models.Card
	if err := json.Unmarshal(buf.Bytes(), &first); err != nil {
		t.Fatalf("decode first card JSON: %v\noutput: %s", err, buf.String())
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "second card",
		"--back", "second back",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var second models.Card
	if err := json.Unmarshal(buf.Bytes(), &second); err != nil {
		t.Fatalf("decode second card JSON: %v\noutput: %s", err, buf.String())
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Deck  string               `json:"deck"`
		Queue []models.QueueEntry  `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode queue JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Deck != "portuguese" {
		t.Fatalf("expected deck portuguese, got %q", resp.Deck)
	}
	if len(resp.Queue) != 2 {
		t.Fatalf("expected 2 queue entries, got %d", len(resp.Queue))
	}
	if resp.Queue[0].Position != 0 || resp.Queue[0].ID != second.ID || resp.Queue[0].FrontPreview != "second card" {
		t.Fatalf("expected second card at position 0, got %+v", resp.Queue[0])
	}
	if resp.Queue[1].Position != 1 || resp.Queue[1].ID != first.ID || resp.Queue[1].FrontPreview != "first card" {
		t.Fatalf("expected first card at position 1, got %+v", resp.Queue[1])
	}
}

func TestQueueEmptyJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	raw := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`"queue": null`)) {
		t.Fatalf("expected empty array, got: %s", raw)
	}

	var resp struct {
		Queue []models.QueueEntry `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode queue JSON: %v\noutput: %s", err, raw)
	}
	if len(resp.Queue) != 0 {
		t.Fatalf("expected 0 queue entries, got %d", len(resp.Queue))
	}
}

func TestQueueDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"queue", "missing", "--json"})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestQueueAfterDelete(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var cardIDs []int64
	for _, spec := range []struct{ front, back string }{
		{"first", "first back"},
		{"second", "second back"},
		{"third", "third back"},
	} {
		buf.Reset()
		rootCmd.SetArgs([]string{
			"add", "portuguese",
			"--front", spec.front,
			"--back", spec.back,
			"--json",
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
		var card models.Card
		if err := json.Unmarshal(buf.Bytes(), &card); err != nil {
			t.Fatalf("decode card JSON: %v\noutput: %s", err, buf.String())
		}
		cardIDs = append(cardIDs, card.ID)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"delete", "portuguese", formatInt(cardIDs[1]), "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Queue []models.QueueEntry `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode queue JSON: %v\noutput: %s", err, buf.String())
	}
	if len(resp.Queue) != 2 {
		t.Fatalf("expected 2 queue entries, got %d", len(resp.Queue))
	}
	if resp.Queue[0].Position != 0 || resp.Queue[0].FrontPreview != "third" {
		t.Fatalf("expected third at position 0, got %+v", resp.Queue[0])
	}
	if resp.Queue[1].Position != 1 || resp.Queue[1].FrontPreview != "first" {
		t.Fatalf("expected first at position 1, got %+v", resp.Queue[1])
	}
}
