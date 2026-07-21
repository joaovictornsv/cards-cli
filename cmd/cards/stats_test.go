package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/queue"
	"github.com/joaovictornsv/cards-cli/internal/study"
)

func TestStatsJSONNeverStudied(t *testing.T) {
	_, buf := testHarness(t)

	rootCmd.SetArgs([]string{"deck", "create", "fresh", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"stats", "fresh", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var stats models.DeckStats
	if err := json.Unmarshal(buf.Bytes(), &stats); err != nil {
		t.Fatalf("decode stats JSON: %v\noutput: %s", err, buf.String())
	}
	if stats.Deck != "fresh" {
		t.Fatalf("deck = %q, want fresh", stats.Deck)
	}
	if stats.SessionsCount != 0 {
		t.Fatalf("sessions_count = %d, want 0", stats.SessionsCount)
	}
	if stats.LastSessionAt != nil {
		t.Fatalf("expected nil last_session_at, got %v", stats.LastSessionAt)
	}
	if stats.LastSessionAgo != "never" {
		t.Fatalf("last_session_ago = %q, want never", stats.LastSessionAgo)
	}
	if stats.Nudge == "" {
		t.Fatal("expected nudge for never studied deck")
	}
}

func TestStatsTableOutput(t *testing.T) {
	_, buf := testHarness(t)

	rootCmd.SetArgs([]string{"deck", "create", "tabledeck"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	jsonOutput = false
	if f := rootCmd.PersistentFlags().Lookup("json"); f != nil {
		_ = f.Value.Set("false")
		f.Changed = false
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"stats", "tabledeck"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{"deck: tabledeck", "sessions: 0", "last session: never", "nudge:"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestStatsDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)

	rootCmd.SetArgs([]string{"stats", "missing", "--json"})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestStudyUpdatesDeckStats(t *testing.T) {
	_, buf := testHarness(t)

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{queue.GradeEasy})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "statsdeck", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"add", "statsdeck", "--front", "hello", "--back", "world", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"study", "statsdeck"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"stats", "statsdeck", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var stats models.DeckStats
	if err := json.Unmarshal(buf.Bytes(), &stats); err != nil {
		t.Fatalf("decode stats JSON: %v\noutput: %s", err, buf.String())
	}
	if stats.SessionsCount != 1 {
		t.Fatalf("sessions_count = %d, want 1", stats.SessionsCount)
	}
	if stats.LastSessionAt == nil || *stats.LastSessionAt == "" {
		t.Fatal("expected last_session_at after study session")
	}
}

func TestStatsStaleSessionNudge(t *testing.T) {
	dbPath, buf := testHarness(t)
	_ = dbPath

	rootCmd.SetArgs([]string{"deck", "create", "stale", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	repo, cleanup, err := openRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	deck, err := repo.GetDeckByName(context.Background(), "stale")
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.RecordDeckSession(context.Background(), deck.ID, "2026-07-01T12:00:00Z"); err != nil {
		t.Fatal(err)
	}
	cleanup()

	buf.Reset()
	rootCmd.SetArgs([]string{"stats", "stale", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var stats models.DeckStats
	if err := json.Unmarshal(buf.Bytes(), &stats); err != nil {
		t.Fatalf("decode stats JSON: %v\noutput: %s", err, buf.String())
	}
	if stats.Nudge == "" {
		t.Fatal("expected nudge for stale session")
	}
	if !strings.Contains(stats.Nudge, "ready for a quick review?") {
		t.Fatalf("unexpected nudge: %q", stats.Nudge)
	}
}

func TestShouldRecordSession(t *testing.T) {
	tests := []struct {
		name   string
		result study.Result
		want   bool
	}{
		{
			name:   "complete",
			result: study.Result{Status: "complete", Reviews: []study.Review{{}}},
			want:   true,
		},
		{
			name:   "quit with reviews",
			result: study.Result{Status: "quit", Reviews: []study.Review{{}}},
			want:   true,
		},
		{
			name:   "quit without reviews",
			result: study.Result{Status: "quit", Reviews: nil},
			want:   false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := shouldRecordSession(tc.result); got != tc.want {
				t.Fatalf("shouldRecordSession() = %v, want %v", got, tc.want)
			}
		})
	}
}
