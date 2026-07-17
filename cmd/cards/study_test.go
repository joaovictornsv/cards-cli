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
			queue.GradeAgain,
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

	wantLabels := []string{"E", "F", "C", "B", "G", "H", "A", "D"}
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

func TestStudyEmptyDeck(t *testing.T) {
	_, _ = testHarness(t)

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput(nil)
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "empty", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"study", "empty"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for empty deck")
	}
	if !strings.Contains(err.Error(), `deck "empty" has no cards`) {
		t.Fatalf("expected friendly empty deck error, got %v", err)
	}
	if !strings.Contains(err.Error(), `cards add empty`) {
		t.Fatalf("expected add-card hint, got %v", err)
	}
}

func TestStudyLimitOverride(t *testing.T) {
	_, buf := testHarness(t)

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{
			queue.GradeEasy,
			queue.GradeEasy,
		})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "limit", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	for _, label := range []string{"H", "G", "F", "E", "D", "C", "B", "A"} {
		rootCmd.SetArgs([]string{
			"add", "limit",
			"--front", label,
			"--back", label,
			"--json",
		})
		if err := rootCmd.Execute(); err != nil {
			t.Fatal(err)
		}
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"study", "limit", "--limit", "2"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "[1/2]") || !strings.Contains(out, "[2/2]") {
		t.Fatalf("expected 2-card batch, got:\n%s", out)
	}
	if strings.Contains(out, "[3/2]") {
		t.Fatalf("expected only 2 cards reviewed, got:\n%s", out)
	}
}

func TestStudyJSONLog(t *testing.T) {
	_, buf := testHarness(t)

	oldJSON := jsonOutput
	jsonOutput = true
	t.Cleanup(func() { jsonOutput = oldJSON })

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{queue.GradeEasy})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "jsonlog", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"add", "jsonlog", "--front", "hello", "--back", "world", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"study", "jsonlog", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Session complete.") {
		t.Fatalf("expected interactive output, got:\n%s", out)
	}

	jsonStart := strings.Index(out, "{\n  \"deck\"")
	if jsonStart < 0 {
		t.Fatalf("expected JSON log in output, got:\n%s", out)
	}

	var result study.Result
	if err := json.Unmarshal([]byte(out[jsonStart:]), &result); err != nil {
		t.Fatalf("decode study JSON log: %v\njson: %s", err, out[jsonStart:])
	}
	if result.Deck != "jsonlog" {
		t.Fatalf("deck = %q, want jsonlog", result.Deck)
	}
	if result.Status != "complete" {
		t.Fatalf("status = %q, want complete", result.Status)
	}
	if len(result.Reviews) != 1 {
		t.Fatalf("reviews = %d, want 1", len(result.Reviews))
	}
	if result.Reviews[0].Grade != queue.GradeEasy {
		t.Fatalf("grade = %q, want easy", result.Reviews[0].Grade)
	}
}

func TestStudyJSONLogReplaceGrade(t *testing.T) {
	_, buf := testHarness(t)

	oldJSON := jsonOutput
	jsonOutput = true
	t.Cleanup(func() { jsonOutput = oldJSON })

	oldFactory := studyInputFactory
	studyInputFactory = func(io.Reader) study.Input {
		return study.NewScriptedInput([]queue.Grade{queue.GradeReplace})
	}
	t.Cleanup(func() { studyInputFactory = oldFactory })

	rootCmd.SetArgs([]string{"deck", "create", "replacejson", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"add", "replacejson", "--front", "hello", "--back", "world", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"study", "replacejson", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	jsonStart := strings.Index(out, "{\n  \"deck\"")
	if jsonStart < 0 {
		t.Fatalf("expected JSON log in output, got:\n%s", out)
	}

	var result study.Result
	if err := json.Unmarshal([]byte(out[jsonStart:]), &result); err != nil {
		t.Fatalf("decode study JSON log: %v\njson: %s", err, out[jsonStart:])
	}
	if len(result.Reviews) != 1 {
		t.Fatalf("reviews = %d, want 1", len(result.Reviews))
	}
	if result.Reviews[0].Grade != queue.GradeReplace {
		t.Fatalf("grade = %q, want replace", result.Reviews[0].Grade)
	}
}

func TestStudyInvalidLimit(t *testing.T) {
	_, _ = testHarness(t)

	rootCmd.SetArgs([]string{"deck", "create", "badlimit", "--json"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"study", "badlimit", "--limit", "0"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid limit")
	}
	if !strings.Contains(err.Error(), "--limit must be at least 1") {
		t.Fatalf("expected limit validation error, got %v", err)
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

	cfg := config.Config{BatchSize: 4, AgainOffset: 2}
	if err := runStudyWithRepo(context.Background(), repo, "small", cfg, 4, nil, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "[1/1]") {
		t.Fatalf("expected [1/1], got:\n%s", out.String())
	}
}
