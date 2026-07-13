package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/config"
	"github.com/joaovictornsv/cards-cli/internal/models"
	"github.com/joaovictornsv/cards-cli/internal/queue"
	"github.com/joaovictornsv/cards-cli/internal/study"
)

func TestStudySessionQueueMutation(t *testing.T) {
	dbPath, buf := testHarness(t)

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{
			queue.GradeEasy,
			queue.GradeAgain,
			queue.GradeHard,
			queue.GradeEasy,
		})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "walk", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	labels := []string{"H", "G", "F", "E", "D", "C", "B", "A"}
	idByLabel := make(map[string]int64)
	for _, label := range labels {
		buf.Reset()
		rootCmd.SetArgs([]string{
			"add", "walk",
			"--front", label,
			"--back", label + " back",
			"--json",
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
		var card models.Card
		if err := json.Unmarshal(buf.Bytes(), &card); err != nil {
			t.Fatal(err)
		}
		idByLabel[label] = card.ID
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"study", "walk"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Session complete.") {
		t.Fatalf("expected session complete, got:\n%s", buf.String())
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"queue", "walk", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	var resp struct {
		Queue []models.QueueEntry `json:"queue"`
	}
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		t.Fatalf("decode queue JSON: %v\noutput: %s", err, buf.String())
	}

	wantLabels := []string{"E", "F", "B", "G", "H", "C", "A", "D"}
	labelByID := make(map[int64]string)
	for label, id := range idByLabel {
		labelByID[id] = label
	}
	if len(resp.Queue) != len(wantLabels) {
		t.Fatalf("queue length = %d, want %d", len(resp.Queue), len(wantLabels))
	}
	for i, entry := range resp.Queue {
		got := labelByID[entry.ID]
		if got != wantLabels[i] {
			t.Fatalf("position %d: got %q, want %q", i, got, wantLabels[i])
		}
	}

	_ = dbPath
}

func TestStudyDeckNotFound(t *testing.T) {
	_, _ = testHarness(t)

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput(nil)
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"study", "missing"})
	err := rootCmd.Execute()
	if !errors.Is(err, errDeckNotFound) {
		t.Fatalf("expected errDeckNotFound, got %v", err)
	}
}

func TestRunStudyWithRepo(t *testing.T) {
	dbPath, _ := testHarness(t)
	_ = dbPath

	rootCmd.SetArgs([]string{"deck", "create", "small", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"add", "small", "--front", "only", "--back", "one", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	repo, cleanup, err := openRepo()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	var out bytes.Buffer
	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{queue.GradeEasy})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	cfg := config.Config{BatchSize: 4, AgainOffset: 2, HardOffset: 5}
	if err := runStudyWithRepo(context.Background(), repo, "small", cfg, nil, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "[1/1]") {
		t.Fatalf("expected [1/1], got:\n%s", out.String())
	}
}
