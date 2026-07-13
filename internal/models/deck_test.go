package models

import (
	"errors"
	"testing"
)

func TestDeckValidateForCreate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantName string
		wantErr  error
	}{
		{name: "valid", input: "portuguese", wantName: "portuguese"},
		{name: "trimmed valid", input: "  portuguese  ", wantName: "portuguese"},
		{name: "whitespace only", input: "   ", wantErr: ErrDeckNameRequired},
		{name: "empty", input: "", wantErr: ErrDeckNameRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := Deck{Name: tt.input}
			err := deck.ValidateForCreate()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if deck.Name != tt.wantName {
				t.Fatalf("expected name %q, got %q", tt.wantName, deck.Name)
			}
		})
	}
}
