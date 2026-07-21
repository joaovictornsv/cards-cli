package models

type DeckStats struct {
	Deck           string  `json:"deck"`
	SessionsCount  int     `json:"sessions_count"`
	LastSessionAt  *string `json:"last_session_at"`
	LastSessionAgo string  `json:"last_session_ago"`
	Nudge          string  `json:"nudge"`
}
