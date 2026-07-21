package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenMemory(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	for _, table := range []string{"decks", "cards", "queue", "schema_migrations"} {
		var name string
		err := database.SQL().QueryRow(
			`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
			table,
		).Scan(&name)
		if err != nil {
			t.Fatalf("table %q: %v", table, err)
		}
	}
}

func TestMigrationsIdempotent(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if err := migrate(database.SQL()); err != nil {
		t.Fatal(err)
	}
}

func TestMigrationReplaceEligibleColumn(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	var name string
	err = database.SQL().QueryRow(`
		SELECT name FROM pragma_table_info('cards') WHERE name = 'replace_eligible'`).Scan(&name)
	if err != nil {
		t.Fatalf("replace_eligible column: %v", err)
	}
}

func TestMigrationDeckStatsColumns(t *testing.T) {
	database, err := OpenMemory()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	for _, col := range []string{"sessions_count", "last_session_at"} {
		var name string
		err = database.SQL().QueryRow(`
			SELECT name FROM pragma_table_info('decks') WHERE name = ?`, col).Scan(&name)
		if err != nil {
			t.Fatalf("%s column: %v", col, err)
		}
	}
}

func TestOpenCreatesParentDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "cards.db")

	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected parent directory to exist: %v", err)
	}
}
