# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `internal/queue` package with pure re-insert algorithm (`again`/`hard`/`easy` grading, batch pull, insert)
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
