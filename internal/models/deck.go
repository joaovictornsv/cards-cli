package models

import (
	"fmt"
	"strings"
	"time"
)

type Deck struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CardCount int    `json:"card_count"`
	CreatedAt string `json:"created_at"`
}

func (d *Deck) ValidateForCreate() error {
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("deck name is required")
	}
	return nil
}

func NowTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
