package models

import (
	"errors"
	"strings"
)

var ErrCardFrontRequired = errors.New("card front is required")

var ErrCardBackRequired = errors.New("card back is required")

var ErrCardEditRequiresField = errors.New("edit requires at least one of --front or --back")

type Card struct {
	ID        int64  `json:"id"`
	DeckID    int64  `json:"deck_id"`
	Front     string `json:"front"`
	Back      string `json:"back"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CardSummary struct {
	ID        int64  `json:"id"`
	Front     string `json:"front"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *Card) NormalizeForCreate() {
	c.Front = strings.TrimSpace(c.Front)
	c.Back = strings.TrimSpace(c.Back)
}

func (c *Card) ValidateForCreate() error {
	c.NormalizeForCreate()
	if c.Front == "" {
		return ErrCardFrontRequired
	}
	if c.Back == "" {
		return ErrCardBackRequired
	}
	return nil
}

// ValidateForUpdate trims and validates partial front/back updates.
// At least one field must be provided; provided fields must be non-empty after trim.
func (c *Card) ValidateForUpdate(front, back *string) error {
	if front == nil && back == nil {
		return ErrCardEditRequiresField
	}
	if front != nil {
		trimmed := strings.TrimSpace(*front)
		if trimmed == "" {
			return ErrCardFrontRequired
		}
		*front = trimmed
	}
	if back != nil {
		trimmed := strings.TrimSpace(*back)
		if trimmed == "" {
			return ErrCardBackRequired
		}
		*back = trimmed
	}
	return nil
}
