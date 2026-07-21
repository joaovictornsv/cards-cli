# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Docs audit: remove duplicate SKILL entry, drop redundant "Available" status columns, update NEXT_STEPS and PROJECT_DRAFT for shipped features

## [0.1.7] - 2026-07-21

### Added

- `cards deck shuffle <name>` — randomly reshuffle the entire deck queue order (`--yes` required with `--json`; no-op for decks with 0–1 cards)

## [0.1.6] - 2026-07-21

### Added

- `cards stats <deck>` — per-deck session count, last-session timestamp, and nudge when study is overdue (`nudge_threshold_days` config, default 3)
- Deck-level `sessions_count` and `last_session_at` updated when a study session completes or is quit after at least one review (migration `004_deck_stats.sql`)

## [0.1.5] - 2026-07-21

### Added

- `cards export <deck>` — export deck cards to JSON or CSV (`--format`, `--output`; `--json` prints summary)
- `cards import` — import cards from JSON or CSV (`--deck`, `--format`, `--file`; `--append` for existing decks)

## [0.1.4] - 2026-07-20

### Added

- `cards search` — find cards across decks by text (OR-matched `--term` flags or positional query; optional `--deck` filter)

## [0.1.3] - 2026-07-17

### Fixed

- Queue compaction after card delete no longer fails with a UNIQUE constraint error on large decks

## [0.1.2] - 2026-07-16

### Added

- Study `replace` grade (`r` / `R`) — same queue behavior as `easy`, sets persistent `replace_eligible` flag on the card
- `replace_eligible` boolean column on cards (migration `003_replace_eligible.sql`)
- `cards list <deck> --replace-eligible` filter for flagged cards
- `cards edit <deck> <id> --replace-eligible=false` to clear the flag
- `replace_eligible` field in `cards list`, `cards show`, and `cards edit` JSON output

## [0.1.1] - 2026-07-16

### Removed

- Study session `hard` grade and `hard_offset` config key — grading is now binary (`again` / `easy`)

## [0.1.0] - 2026-07-16

### Added

- Migration `002_drop_card_stats.sql` — remove unused per-card stats columns from schema
- Docs audit: concise agent guidance, fixed `deck delete --yes` with `--json`, removed deferred `stats` references
- `cards study --limit` batch size override (default from config `batch_size`)
- `cards study --json` session log emitted after interactive output (deck, batch size, reviews, status)
- Friendly empty-deck error with hint to use `cards add`
- `internal/study` `Result` and `Review` types for session logging
- `PrintStudyLog` formatter in `internal/output`
- `internal/study` package with testable session engine and terminal input handling
- `ListQueueCardIDsByDeck` and `ReplaceDeckQueue` repository methods for study queue updates
- `cards queue` command with `--json` output for inspecting deck queue order
- `ListQueueByDeck` repository method and `QueueEntry` model
- Queue formatters in `internal/output` (table and JSON)
- `cards show`, `cards edit`, and `cards delete` commands with `--json` output
- Card repository methods `GetCardByDeckAndID`, `UpdateCard`, and `DeleteCard` (with queue compaction on delete)
- `ValidateForUpdate` partial-update validation in `internal/models`
- `cards add` and `cards list` commands with `--json` output
- `Card` and `CardSummary` types with validation in `internal/models`
- Card repository methods in `internal/db` (`CreateCard`, `ListCardsByDeck`, `GetCardByID`)
- Card formatters in `internal/output` (table and JSON)
- `cards deck create`, `cards deck list`, and `cards deck delete` commands with `--json` output
- `internal/models` package with `Deck` type and validation
- Deck repository methods in `internal/db` (`CreateDeck`, `GetDeckByName`, `ListDecks`, `DeleteDeckByName`)
- Deck formatters in `internal/output` (table and JSON)
- SQLite foreign key enforcement (`PRAGMA foreign_keys = ON`) for cascade deletes
- SQLite schema (`decks`, `cards`, `queue`) with embedded migrations
- `internal/db` package with `Open` / `OpenMemory` and idempotent migration runner (`modernc.org/sqlite`)
- `internal/output` package with table and JSON formatters for `config` and `version`
- `runWithRepo` helper in CLI root for future DB-backed commands

## [0.0.0] - TBD

Initial v1 release (planned).
