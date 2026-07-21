package stats

import (
	"fmt"
	"time"

	"github.com/joaovictornsv/cards-cli/internal/models"
)

const DefaultNudgeThresholdDays = 3

func BuildDeckStats(deck string, sessionsCount int, lastSessionAt *string, thresholdDays int, now time.Time) models.DeckStats {
	if thresholdDays < 1 {
		thresholdDays = DefaultNudgeThresholdDays
	}

	stats := models.DeckStats{
		Deck:          deck,
		SessionsCount: sessionsCount,
		LastSessionAt: lastSessionAt,
	}

	if lastSessionAt == nil || *lastSessionAt == "" {
		stats.LastSessionAgo = "never"
		stats.Nudge = "never studied — ready for a quick review?"
		return stats
	}

	parsed, err := time.Parse(time.RFC3339, *lastSessionAt)
	if err != nil {
		stats.LastSessionAgo = *lastSessionAt
		return stats
	}

	ago := formatRelativeTime(parsed, now)
	stats.LastSessionAgo = ago

	threshold := time.Duration(thresholdDays) * 24 * time.Hour
	if now.Sub(parsed) >= threshold {
		stats.Nudge = fmt.Sprintf("last session: %s — ready for a quick review?", ago)
	}

	return stats
}

func formatRelativeTime(t time.Time, now time.Time) string {
	t = t.UTC()
	now = now.UTC()

	tDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	days := int(nowDay.Sub(tDay).Hours() / 24)

	switch {
	case days == 0:
		return "today"
	case days == 1:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", days)
	}
}
