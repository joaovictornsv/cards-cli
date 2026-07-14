package study

import "github.com/joaovictornsv/cards-cli/internal/queue"

type Review struct {
	CardID   int64       `json:"card_id"`
	Front    string      `json:"front"`
	Grade    queue.Grade `json:"grade"`
	Position int         `json:"position"`
}

type Result struct {
	Deck      string   `json:"deck"`
	BatchSize int      `json:"batch_size"`
	DeckSize  int      `json:"deck_size"`
	Status    string   `json:"status"`
	Reviews   []Review `json:"reviews"`
}
