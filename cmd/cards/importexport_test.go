package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/importexport"
	"github.com/joaovictornsv/cards-cli/internal/models"
)

func createDeckWithCards(t *testing.T) {
	t.Helper()
	rootCmd.SetArgs([]string{"deck", "create", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	for _, card := range []struct{ front, back string }{
		{"first", "first back"},
		{"second", "second back"},
		{"third", "third back"},
	} {
		resetCommandFlags(t)
		rootCmd.SetArgs([]string{
			"add", "portuguese",
			"--front", card.front,
			"--back", card.back,
			"--json",
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestExportJSON(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithCards(t)

	resetCommandFlags(t)
	buf.Reset()
	rootCmd.SetArgs([]string{"export", "portuguese", "--format", "json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var data importexport.DeckExport
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("decode export JSON: %v\noutput: %s", err, buf.String())
	}
	if data.Deck != "portuguese" {
		t.Fatalf("expected deck portuguese, got %q", data.Deck)
	}
	if len(data.Cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(data.Cards))
	}
	if data.Cards[0].Front != "third" {
		t.Fatalf("expected queue order (third first), got %q", data.Cards[0].Front)
	}
}

func TestExportCSV(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithCards(t)

	resetCommandFlags(t)
	buf.Reset()
	rootCmd.SetArgs([]string{"export", "portuguese", "--format", "csv"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "front,back\n") {
		t.Fatalf("expected CSV header, got: %s", output)
	}
	if !strings.Contains(output, "third,third back") {
		t.Fatalf("expected third card in CSV, got: %s", output)
	}
}

func TestExportJSONSummary(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithCards(t)

	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "export.json")

	buf.Reset()
	rootCmd.SetArgs([]string{
		"export", "portuguese",
		"--format", "json",
		"--output", outPath,
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var summary importexport.ExportSummary
	if err := json.Unmarshal(buf.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v\noutput: %s", err, buf.String())
	}
	if summary.Deck != "portuguese" || summary.Format != "json" || summary.CardCount != 3 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if summary.Output != outPath {
		t.Fatalf("expected output path %q, got %q", outPath, summary.Output)
	}

	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	var data importexport.DeckExport
	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Cards) != 3 {
		t.Fatalf("expected 3 cards in file, got %d", len(data.Cards))
	}
}

func TestExportDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)
	rootCmd.SetArgs([]string{"export", "missing", "--format", "json"})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestImportJSONCreatesDeck(t *testing.T) {
	_, buf := testHarness(t)
	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "deck.json")
	content := `{
  "deck": "spanish",
  "cards": [
    {"front": "hola", "back": "hello"},
    {"front": "adiós", "back": "goodbye"}
  ]
}`
	if err := os.WriteFile(importPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"import",
		"--deck", "spanish",
		"--format", "json",
		"--file", importPath,
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result importexport.ImportResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("decode import result: %v\noutput: %s", err, buf.String())
	}
	if result.CardsImported != 2 {
		t.Fatalf("expected 2 imported, got %d", result.CardsImported)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "spanish", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var queueResp struct {
		Queue []models.QueueEntry `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &queueResp); err != nil {
		t.Fatal(err)
	}
	if len(queueResp.Queue) != 2 {
		t.Fatalf("expected 2 queue entries, got %d", len(queueResp.Queue))
	}
	if queueResp.Queue[0].FrontPreview != "hola" {
		t.Fatalf("expected hola at front, got %q", queueResp.Queue[0].FrontPreview)
	}
	if queueResp.Queue[1].FrontPreview != "adiós" {
		t.Fatalf("expected adiós second, got %q", queueResp.Queue[1].FrontPreview)
	}
}

func TestImportCSV(t *testing.T) {
	_, buf := testHarness(t)
	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "deck.csv")
	content := "front,back\nhola,hello\n"
	if err := os.WriteFile(importPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"import",
		"--deck", "spanish",
		"--format", "csv",
		"--file", importPath,
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result importexport.ImportResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.CardsImported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.CardsImported)
	}
}

func TestImportAppend(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithCards(t)

	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "append.json")
	content := `{"deck":"portuguese","cards":[{"front":"new","back":"new back"}]}`
	if err := os.WriteFile(importPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"import",
		"--deck", "portuguese",
		"--format", "json",
		"--file", importPath,
		"--append",
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "portuguese", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var queueResp struct {
		Queue []models.QueueEntry `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &queueResp); err != nil {
		t.Fatal(err)
	}
	if len(queueResp.Queue) != 4 {
		t.Fatalf("expected 4 queue entries, got %d", len(queueResp.Queue))
	}
	if queueResp.Queue[0].FrontPreview != "new" {
		t.Fatalf("expected new card at front, got %q", queueResp.Queue[0].FrontPreview)
	}
}

func TestImportDeckAlreadyExists(t *testing.T) {
	_, _ = testHarness(t)
	createDeckWithCards(t)

	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "deck.json")
	content := `{"deck":"portuguese","cards":[{"front":"x","back":"y"}]}`
	if err := os.WriteFile(importPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{
		"import",
		"--deck", "portuguese",
		"--format", "json",
		"--file", importPath,
	})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeckAlreadyExists) {
		t.Fatalf("expected errDeckAlreadyExists, got %v", err)
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	_, buf := testHarness(t)
	createDeckWithCards(t)

	tmp := t.TempDir()
	exportPath := filepath.Join(tmp, "export.json")
	importPath := filepath.Join(tmp, "import-target.json")

	buf.Reset()
	rootCmd.SetArgs([]string{"export", "portuguese", "--format", "json", "--output", exportPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"deck", "delete", "portuguese", "--json", "--yes"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatal(err)
	}
	rewritten := strings.Replace(string(raw), `"portuguese"`, `"roundtrip"`, 1)
	if err := os.WriteFile(importPath, []byte(rewritten), 0o644); err != nil {
		t.Fatal(err)
	}

	resetCommandFlags(t)
	buf.Reset()
	rootCmd.SetArgs([]string{
		"import",
		"--deck", "roundtrip",
		"--format", "json",
		"--file", importPath,
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	resetCommandFlags(t)
	buf.Reset()
	rootCmd.SetArgs([]string{"export", "roundtrip", "--format", "json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var reimported importexport.DeckExport
	if err := json.Unmarshal(buf.Bytes(), &reimported); err != nil {
		t.Fatal(err)
	}
	if len(reimported.Cards) != 3 {
		t.Fatalf("expected 3 cards after round trip, got %d", len(reimported.Cards))
	}
	if reimported.Cards[0].Front != "third" || reimported.Cards[2].Front != "first" {
		t.Fatalf("queue order not preserved: %+v", reimported.Cards)
	}
}

func TestImportMalformedRowsReported(t *testing.T) {
	_, buf := testHarness(t)
	tmp := t.TempDir()
	importPath := filepath.Join(tmp, "partial.csv")
	content := "front,back\nvalid,valid\nonly-one\n"
	if err := os.WriteFile(importPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{
		"import",
		"--deck", "partial",
		"--format", "csv",
		"--file", importPath,
		"--json",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var result importexport.ImportResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.CardsImported != 1 {
		t.Fatalf("expected 1 imported, got %d", result.CardsImported)
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 parse error, got %d: %v", len(result.Errors), result.Errors)
	}
}
