package importexport

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestWriteAndParseJSONRoundTrip(t *testing.T) {
	id := int64(42)
	original := DeckExport{
		Deck: "portuguese",
		Cards: []CardExport{
			{ID: &id, Front: "What is saudade?", Back: "A deep emotional state."},
			{Front: "Olá", Back: "Hello"},
		},
	}

	var buf bytes.Buffer
	if err := WriteJSON(&buf, original); err != nil {
		t.Fatal(err)
	}

	parsed, err := ParseJSON(&buf, "portuguese")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Deck != "portuguese" {
		t.Fatalf("expected deck portuguese, got %q", parsed.Deck)
	}
	if len(parsed.Cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(parsed.Cards))
	}
	if parsed.Cards[0].Front != "What is saudade?" || parsed.Cards[0].Back != "A deep emotional state." {
		t.Fatalf("unexpected first card: %+v", parsed.Cards[0])
	}
}

func TestParseJSONDeckMismatch(t *testing.T) {
	input := `{"deck":"other","cards":[{"front":"a","back":"b"}]}`
	_, err := ParseJSON(strings.NewReader(input), "portuguese")
	if !errors.Is(err, ErrDeckNameMismatch) {
		t.Fatalf("expected ErrDeckNameMismatch, got %v", err)
	}
}

func TestParseJSONInvalid(t *testing.T) {
	_, err := ParseJSON(strings.NewReader(`{invalid`), "portuguese")
	if !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestWriteAndParseCSVRoundTrip(t *testing.T) {
	cards := []CardExport{
		{Front: "What is saudade?", Back: "A deep emotional state."},
		{Front: "Olá", Back: "Hello"},
	}

	var buf bytes.Buffer
	if err := WriteCSV(&buf, cards); err != nil {
		t.Fatal(err)
	}

	parsed, errs, err := ParseCSV(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(parsed) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(parsed))
	}
	if parsed[0].Front != "What is saudade?" {
		t.Fatalf("unexpected first card front: %q", parsed[0].Front)
	}
}

func TestParseCSVEmptyFile(t *testing.T) {
	parsed, errs, err := ParseCSV(strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 0 || len(errs) != 0 {
		t.Fatalf("expected empty result, got cards=%d errs=%d", len(parsed), len(errs))
	}
}

func TestParseCSVHeaderOnly(t *testing.T) {
	parsed, errs, err := ParseCSV(strings.NewReader("front,back\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 0 || len(errs) != 0 {
		t.Fatalf("expected empty result, got cards=%d errs=%d", len(parsed), len(errs))
	}
}

func TestParseCSVMalformedRow(t *testing.T) {
	input := "front,back\nvalid front,valid back\nonly-one-field\n,both empty\n"
	parsed, errs, err := ParseCSV(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 1 {
		t.Fatalf("expected 1 valid card, got %d", len(parsed))
	}
	if len(errs) != 2 {
		t.Fatalf("expected 2 row errors, got %d: %v", len(errs), errs)
	}
}

func TestParseCSVInvalidHeader(t *testing.T) {
	_, _, err := ParseCSV(strings.NewReader("question,answer\na,b\n"))
	if !errors.Is(err, ErrInvalidCSVHeader) {
		t.Fatalf("expected ErrInvalidCSVHeader, got %v", err)
	}
}

func TestParseCSVSkipsBlankLines(t *testing.T) {
	input := "front,back\n\n  a  ,  b  \n\n"
	parsed, errs, err := ParseCSV(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(parsed) != 1 || parsed[0].Front != "a" || parsed[0].Back != "b" {
		t.Fatalf("unexpected parsed card: %+v", parsed)
	}
}

func TestCardsFromExportValidation(t *testing.T) {
	data := DeckExport{
		Deck: "portuguese",
		Cards: []CardExport{
			{Front: "valid", Back: "valid"},
			{Front: "", Back: "missing front"},
		},
	}
	cards, errs := CardsFromExport(data)
	if len(cards) != 1 {
		t.Fatalf("expected 1 valid card, got %d", len(cards))
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}
