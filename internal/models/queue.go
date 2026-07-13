package models

type QueueEntry struct {
	Position     int    `json:"position"`
	ID           int64  `json:"id"`
	FrontPreview string `json:"front_preview"`
}
