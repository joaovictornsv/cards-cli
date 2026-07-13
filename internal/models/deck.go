package models

import (
	"errors"
	"strings"
	"time"
)

var ErrDeckNameRequired = errors.New("deck name is required")

type Deck struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CardCount int    `json:"card_count"`
	CreatedAt string `json:"created_at"`
}

func (d *Deck) NormalizeForCreate() {
	d.Name = strings.TrimSpace(d.Name)
}

func (d *Deck) ValidateForCreate() error {
	d.NormalizeForCreate()
	if d.Name == "" {
		return ErrDeckNameRequired
	}
	return nil
}

func NowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
