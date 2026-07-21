package importexport

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

func WriteCSV(w io.Writer, cards []CardExport) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"front", "back"}); err != nil {
		return err
	}
	for _, card := range cards {
		if err := writer.Write([]string{card.Front, card.Back}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func ParseCSV(r io.Reader) ([]CardInput, []string, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	if len(records) == 0 {
		return []CardInput{}, nil, nil
	}

	start := 0
	header := records[0]
	if len(header) >= 2 && strings.EqualFold(strings.TrimSpace(header[0]), "front") &&
		strings.EqualFold(strings.TrimSpace(header[1]), "back") {
		start = 1
	} else if !isBlankRow(header) {
		return nil, nil, ErrInvalidCSVHeader
	}

	cards := make([]CardInput, 0)
	errs := make([]string, 0)
	for i := start; i < len(records); i++ {
		row := records[i]
		if isBlankRow(row) {
			continue
		}
		if len(row) < 2 {
			errs = append(errs, fmt.Sprintf("row %d: %v", i+1, ErrMissingCSVField))
			continue
		}
		input, err := ValidateCardInput(row[0], row[1])
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %v", i+1, err))
			continue
		}
		cards = append(cards, input)
	}
	return cards, errs, nil
}

func isBlankRow(row []string) bool {
	for _, field := range row {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}
