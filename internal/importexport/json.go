package importexport

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteJSON(w io.Writer, data DeckExport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func ParseJSON(r io.Reader, expectedDeck string) (DeckExport, error) {
	var data DeckExport
	dec := json.NewDecoder(r)
	if err := dec.Decode(&data); err != nil {
		return DeckExport{}, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if data.Deck != "" && expectedDeck != "" && data.Deck != expectedDeck {
		return DeckExport{}, ErrDeckNameMismatch
	}
	if data.Deck == "" {
		data.Deck = expectedDeck
	}
	return data, nil
}

func CardsFromExport(data DeckExport) ([]CardInput, []string) {
	cards := make([]CardInput, 0, len(data.Cards))
	errs := make([]string, 0)
	for i, card := range data.Cards {
		input, err := ValidateCardInput(card.Front, card.Back)
		if err != nil {
			errs = append(errs, fmt.Sprintf("card %d: %v", i+1, err))
			continue
		}
		cards = append(cards, input)
	}
	return cards, errs
}
