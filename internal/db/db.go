package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql *sql.DB
}

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	sqlDB, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := enableForeignKeys(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	if err := migrate(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &DB{sql: sqlDB}, nil
}

func OpenMemory() (*DB, error) {
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open in-memory database: %w", err)
	}

	if err := enableForeignKeys(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	if err := migrate(sqlDB); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return &DB{sql: sqlDB}, nil
}

func enableForeignKeys(sqlDB *sql.DB) error {
	if _, err := sqlDB.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}
	return nil
}

func (d *DB) Close() error {
	return d.sql.Close()
}

func (d *DB) SQL() *sql.DB {
	return d.sql
}
