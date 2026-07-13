package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/db"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestAddJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
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

	var card models.Card
	if err := json.Unmarshal(buf.Bytes(), &card); err != nil {
		t.Fatalf("decode card JSON: %v\noutput: %s", err, buf.String())
	}
	if card.ID <= 0 {
		t.Fatalf("expected positive id, got %d", card.ID)
	}
	if card.Front != "What is saudade?" {
		t.Fatalf("expected front %q, got %q", "What is saudade?", card.Front)
	}
	if card.Back != "A deep emotional state of longing." {
		t.Fatalf("expected back %q, got %q", "A deep emotional state of longing.", card.Back)
	}
}

func TestAddAtQueueFront(t *testing.T) {
	dbPath, _ := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "first",
		"--back", "first back",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "second",
		"--back", "second back",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	var firstPos, secondPos int
	err = database.SQL().QueryRow(`
		SELECT q.position
		FROM queue q
		JOIN cards c ON c.id = q.card_id
		WHERE c.front = ?
		ORDER BY q.position`,
		"first",
	).Scan(&firstPos)
	if errors.Is(err, sql.ErrNoRows) {
		t.Fatal("expected first card in queue")
	}
	if err != nil {
		t.Fatal(err)
	}

	err = database.SQL().QueryRow(`
		SELECT q.position
		FROM queue q
		JOIN cards c ON c.id = q.card_id
		WHERE c.front = ?`,
		"second",
	).Scan(&secondPos)
	if err != nil {
		t.Fatal(err)
	}
	if secondPos != 0 {
		t.Fatalf("expected newest card at position 0, got %d", secondPos)
	}
	if firstPos != 1 {
		t.Fatalf("expected first card at position 1, got %d", firstPos)
	}
}

func TestAddDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{
		"add", "missing",
		"--front", "front",
		"--back", "back",
		"--json",
	})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestAddValidatesFlags(t *testing.T) {
	dbPath, _ := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	resetCommandFlags(t)
	rootCmd.SetArgs([]string{"add", "portuguese", "--back", "back"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for missing --front")
	}

	resetCommandFlags(t)
	rootCmd.SetArgs([]string{"add", "portuguese", "--front", "front"})
	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for missing --back")
	}

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	var count int
	if err := database.SQL().QueryRow(`SELECT COUNT(*) FROM cards`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("expected no cards inserted, got %d", count)
	}
}

func TestAddValidatesBeforeDB(t *testing.T) {
	dbPath, _ := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{
		"add", "portuguese",
		"--front", "   ",
		"--back", "back",
	})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error for whitespace front")
	}
	if !errors.Is(err, models.ErrCardFrontRequired) {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(dbPath); statErr != nil {
		t.Fatal(statErr)
	}
}

func TestListJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

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
	rootCmd.SetArgs([]string{"list", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Deck  string                `json:"deck"`
		Cards []models.CardSummary  `json:"cards"`
		Total int                   `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode list JSON: %v\noutput: %s", err, buf.String())
	}
	if resp.Deck != "portuguese" {
		t.Fatalf("expected deck portuguese, got %q", resp.Deck)
	}
	if len(resp.Cards) != 1 {
		t.Fatalf("expected 1 card, got %d", len(resp.Cards))
	}
	if resp.Total != 1 {
		t.Fatalf("expected total 1, got %d", resp.Total)
	}
	if resp.Cards[0].Front != "What is saudade?" {
		t.Fatalf("expected front %q, got %q", "What is saudade?", resp.Cards[0].Front)
	}
	if resp.Cards[0].CreatedAt == "" || resp.Cards[0].UpdatedAt == "" {
		t.Fatalf("expected timestamps, got created_at=%q updated_at=%q", resp.Cards[0].CreatedAt, resp.Cards[0].UpdatedAt)
	}
}

func TestListEmptyJSON(t *testing.T) {
	_, buf := testHarness(t)
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"list", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	raw := buf.String()
	if bytes.Contains(buf.Bytes(), []byte(`"cards": null`)) {
		t.Fatalf("expected empty array, got: %s", raw)
	}

	var resp struct {
		Cards []models.CardSummary `json:"cards"`
		Total int                  `json:"total"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode list JSON: %v\noutput: %s", err, raw)
	}
	if len(resp.Cards) != 0 {
		t.Fatalf("expected 0 cards, got %d", len(resp.Cards))
	}
	if resp.Total != 0 {
		t.Fatalf("expected total 0, got %d", resp.Total)
	}
}

func TestListDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"list", "missing", "--json"})

	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}
