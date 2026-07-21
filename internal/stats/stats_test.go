package stats

import (
	"testing"
	"time"
)

func TestBuildDeckStatsNeverStudied(t *testing.T) {
	stats := BuildDeckStats("portuguese", 0, nil, 3, time.Now())
	if stats.LastSessionAgo != "never" {
		t.Fatalf("last_session_ago = %q, want never", stats.LastSessionAgo)
	}
	if stats.Nudge == "" {
		t.Fatal("expected nudge for never studied deck")
	}
}

func TestBuildDeckStatsRecentSession(t *testing.T) {
	now := time.Date(2026, 7, 21, 12, 0, 0, 0, time.UTC)
	last := now.Add(-24 * time.Hour).Format(time.RFC3339)
	stats := BuildDeckStats("portuguese", 2, &last, 3, now)
	if stats.LastSessionAgo != "yesterday" {
		t.Fatalf("last_session_ago = %q, want yesterday", stats.LastSessionAgo)
	}
	if stats.Nudge != "" {
		t.Fatalf("expected no nudge, got %q", stats.Nudge)
	}
}

func TestBuildDeckStatsStaleSession(t *testing.T) {
	now := time.Date(2026, 7, 21, 12, 0, 0, 0, time.UTC)
	last := now.Add(-4 * 24 * time.Hour).Format(time.RFC3339)
	stats := BuildDeckStats("portuguese", 5, &last, 3, now)
	if stats.LastSessionAgo != "4 days ago" {
		t.Fatalf("last_session_ago = %q, want 4 days ago", stats.LastSessionAgo)
	}
	if stats.Nudge == "" {
		t.Fatal("expected nudge for stale session")
	}
	if stats.Nudge != "last session: 4 days ago — ready for a quick review?" {
		t.Fatalf("unexpected nudge: %q", stats.Nudge)
	}
}
