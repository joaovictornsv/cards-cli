package models

import (
	"errors"
	"testing"
)

func TestCardValidateForCreate(t *testing.T) {
	tests := []struct {
		name      string
		front     string
		back      string
		wantFront string
		wantBack  string
		wantErr   error
	}{
		{
			name:      "valid",
			front:     "What is saudade?",
			back:      "A deep emotional state of longing.",
			wantFront: "What is saudade?",
			wantBack:  "A deep emotional state of longing.",
		},
		{
			name:      "trimmed valid",
			front:     "  front  ",
			back:      "  back  ",
			wantFront: "front",
			wantBack:  "back",
		},
		{
			name:    "empty front",
			front:   "",
			back:    "back",
			wantErr: ErrCardFrontRequired,
		},
		{
			name:    "whitespace front",
			front:   "   ",
			back:    "back",
			wantErr: ErrCardFrontRequired,
		},
		{
			name:    "empty back",
			front:   "front",
			back:    "",
			wantErr: ErrCardBackRequired,
		},
		{
			name:    "whitespace back",
			front:   "front",
			back:    "   ",
			wantErr: ErrCardBackRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := Card{Front: tt.front, Back: tt.back}
			err := card.ValidateForCreate()
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if card.Front != tt.wantFront {
				t.Fatalf("expected front %q, got %q", tt.wantFront, card.Front)
			}
			if card.Back != tt.wantBack {
				t.Fatalf("expected back %q, got %q", tt.wantBack, card.Back)
			}
		})
	}
}

func TestCardValidateForUpdate(t *testing.T) {
	tests := []struct {
		name    string
		front   *string
		back    *string
		wantErr error
	}{
		{
			name:  "front only",
			front: strPtr("new front"),
		},
		{
			name: "back only",
			back: strPtr("new back"),
		},
		{
			name:  "both",
			front: strPtr("new front"),
			back:  strPtr("new back"),
		},
		{
			name:    "neither",
			wantErr: ErrCardEditRequiresField,
		},
		{
			name:    "empty front",
			front:   strPtr(""),
			wantErr: ErrCardFrontRequired,
		},
		{
			name:    "whitespace front",
			front:   strPtr("   "),
			wantErr: ErrCardFrontRequired,
		},
		{
			name:    "empty back",
			back:    strPtr(""),
			wantErr: ErrCardBackRequired,
		},
		{
			name:    "whitespace back",
			back:    strPtr("   "),
			wantErr: ErrCardBackRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var card Card
			err := card.ValidateForUpdate(tt.front, tt.back)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
