package buildinfo

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()
	if info.Version == "" {
		t.Fatal("expected non-empty version")
	}
	if info.Commit == "" {
		t.Fatal("expected non-empty commit")
	}
	if !strings.HasPrefix(info.GoVersion, "go") {
		t.Fatalf("unexpected go_version: %q", info.GoVersion)
	}
}
