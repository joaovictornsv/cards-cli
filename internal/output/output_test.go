package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/buildinfo"
	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

func TestJSONFormatter(t *testing.T) {
	cfg := config.Config{
		DatabasePath: "/home/user/.local/share/cards/cards.db",
		ConfigPath:   "/home/user/.config/cards/config.toml",
		ConfigExists: false,
		Source:       config.SourceDefault,
		BatchSize:   4,
		AgainOffset: 2,
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
		BatchSize:   4,
		AgainOffset: 2,
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

func TestDeckFormatters(t *testing.T) {
	deck := models.Deck{
		ID:        1,
		Name:      "portuguese",
		CardCount: 0,
		CreatedAt: "2026-07-09T12:00:00Z",
	}

	var buf bytes.Buffer
	jsonFmt := JSONFormatter{}
	if err := jsonFmt.PrintDeck(&buf, deck); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"id": 1`,
		`"name": "portuguese"`,
		`"card_count": 0`,
		`"created_at": "2026-07-09T12:00:00Z"`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in deck json, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintDecks(&buf, []models.Deck{deck}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"decks"`) {
		t.Fatalf("expected decks key in json, got: %s", buf.String())
	}

	buf.Reset()
	table := TableFormatter{}
	if err := table.PrintDecks(&buf, []models.Deck{deck}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ID", "NAME", "CARDS", "portuguese"} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in deck table, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintDecks(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), `"decks": null`) {
		t.Fatalf("expected empty array, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"decks": []`) {
		t.Fatalf("expected empty decks array, got: %s", buf.String())
	}
}

func TestCardFormatters(t *testing.T) {
	card := models.Card{
		ID:              1,
		DeckID:          2,
		Front:           "What is saudade?",
		Back:            "A deep emotional state of longing.",
		CreatedAt:       "2026-07-09T12:00:00Z",
		UpdatedAt:       "2026-07-09T12:00:00Z",
		ReplaceEligible: false,
	}
	summary := models.CardSummary{
		ID:              1,
		Front:           "What is saudade?",
		CreatedAt:       "2026-07-09T12:00:00Z",
		UpdatedAt:       "2026-07-09T12:00:00Z",
		ReplaceEligible: false,
	}

	var buf bytes.Buffer
	jsonFmt := JSONFormatter{}
	if err := jsonFmt.PrintCard(&buf, card); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"id": 1`,
		`"deck_id": 2`,
		`"front": "What is saudade?"`,
		`"back": "A deep emotional state of longing."`,
		`"created_at": "2026-07-09T12:00:00Z"`,
		`"updated_at": "2026-07-09T12:00:00Z"`,
		`"replace_eligible": false`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in card json, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintCards(&buf, "portuguese", []models.CardSummary{summary}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"deck": "portuguese"`,
		`"cards"`,
		`"total": 1`,
		`"front": "What is saudade?"`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in cards json, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintCards(&buf, "portuguese", nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), `"cards": null`) {
		t.Fatalf("expected empty array, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"cards": []`) {
		t.Fatalf("expected empty cards array, got: %s", buf.String())
	}

	buf.Reset()
	table := TableFormatter{}
	if err := table.PrintCards(&buf, "portuguese", []models.CardSummary{summary}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"deck: portuguese", "ID", "FRONT", "CREATED", "UPDATED", "REPLACE", "What is saudade?"} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in cards table, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := table.PrintCard(&buf, card); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"id: 1",
		"front: What is saudade?",
		"back: A deep emotional state of longing.",
		"created_at: 2026-07-09T12:00:00Z",
		"updated_at: 2026-07-09T12:00:00Z",
		"replace_eligible: false",
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in card table, got: %s", want, buf.String())
		}
	}
}

func TestSearchFormatters(t *testing.T) {
	result := models.CardSearchResult{
		ID:    1,
		Deck:  "portuguese",
		Front: "What is saudade?",
		Back:  "A deep emotional state of longing.",
	}

	var buf bytes.Buffer
	jsonFmt := JSONFormatter{}
	if err := jsonFmt.PrintSearchResults(&buf, []models.CardSearchResult{result}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"cards"`,
		`"total": 1`,
		`"deck": "portuguese"`,
		`"front": "What is saudade?"`,
		`"back": "A deep emotional state of longing."`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in search json, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintSearchResults(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), `"cards": null`) {
		t.Fatalf("expected empty array, got: %s", buf.String())
	}

	buf.Reset()
	table := TableFormatter{}
	if err := table.PrintSearchResults(&buf, []models.CardSearchResult{result}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"DECK", "ID", "FRONT", "BACK", "portuguese", "saudade"} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in search table, got: %s", want, buf.String())
		}
	}
}

func TestQueueFormatters(t *testing.T) {
	entry := models.QueueEntry{
		Position:     0,
		ID:           3,
		FrontPreview: "What is saudade?",
	}

	var buf bytes.Buffer
	jsonFmt := JSONFormatter{}
	if err := jsonFmt.PrintQueue(&buf, "portuguese", []models.QueueEntry{entry}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"deck": "portuguese"`,
		`"queue"`,
		`"position": 0`,
		`"id": 3`,
		`"front_preview": "What is saudade?"`,
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in queue json, got: %s", want, buf.String())
		}
	}

	buf.Reset()
	if err := jsonFmt.PrintQueue(&buf, "portuguese", nil); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), `"queue": null`) {
		t.Fatalf("expected empty array, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), `"queue": []`) {
		t.Fatalf("expected empty queue array, got: %s", buf.String())
	}

	buf.Reset()
	table := TableFormatter{}
	if err := table.PrintQueue(&buf, "portuguese", []models.QueueEntry{entry}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"deck: portuguese", "POSITION", "ID", "FRONT", "What is saudade?"} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("expected %q in queue table, got: %s", want, buf.String())
		}
	}
}
