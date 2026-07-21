package models

type ShuffleResult struct {
	Deck      string `json:"deck"`
	CardCount int    `json:"card_count"`
	Status    string `json:"status"`
}
