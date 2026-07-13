package study

import (
	"errors"
	"testing"

	"github.com/joaovictornsv/cards-cli/internal/queue"
)

func TestParseGradeBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    queue.Grade
		wantOk  bool
		wantErr error
	}{
		{name: "1 again", input: []byte{'1'}, want: queue.GradeAgain, wantOk: true},
		{name: "2 hard", input: []byte{'2'}, want: queue.GradeHard, wantOk: true},
		{name: "3 easy", input: []byte{'3'}, want: queue.GradeEasy, wantOk: true},
		{name: "up arrow", input: []byte{27, '[', 'A'}, want: queue.GradeAgain, wantOk: true},
		{name: "down arrow", input: []byte{27, '[', 'B'}, want: queue.GradeHard, wantOk: true},
		{name: "right arrow", input: []byte{27, '[', 'C'}, want: queue.GradeEasy, wantOk: true},
		{name: "q quit", input: []byte{'q'}, wantErr: ErrQuit},
		{name: "unknown", input: []byte{'x'}, wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := parseGradeBytes(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("parseGradeBytes() err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseGradeBytes() unexpected err: %v", err)
			}
			if ok != tt.wantOk {
				t.Fatalf("parseGradeBytes() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Fatalf("parseGradeBytes() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsRevealKey(t *testing.T) {
	for _, key := range []byte{' ', '\r', '\n'} {
		done, err := isRevealKey([]byte{key})
		if err != nil || !done {
			t.Fatalf("isRevealKey(%q) = %v, %v; want true, nil", key, done, err)
		}
	}

	done, err := isRevealKey([]byte{'q'})
	if !errors.Is(err, ErrQuit) || done {
		t.Fatalf("isRevealKey(q) = %v, %v; want false, ErrQuit", done, err)
	}
}
