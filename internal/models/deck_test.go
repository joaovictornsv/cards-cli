package models

import "testing"

func TestDeckValidateForCreate(t *testing.T) {
	deck := Deck{Name: "portuguese"}
	if err := deck.ValidateForCreate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	deck.Name = "   "
	if err := deck.ValidateForCreate(); err == nil {
		t.Fatal("expected error for empty name")
	}
}
