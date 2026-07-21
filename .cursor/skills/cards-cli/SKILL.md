---
name: cards-cli
description: >-
  Manage flashcard decks and cards via the cards CLI (deck create/list/delete,
  add, list, search, show, edit, delete, queue, export, import, config). Use when the user
  mentions flashcards, decks, cards, study queue, spaced repetition queue,
  or cards-cli. Do NOT use for interactive study sessions — those are user-run.
---

# cards-cli Agent

Operate the `cards` CLI in the shell for **deck and card management** — never simulate database changes.

**Do not** explore `cmd/`, `internal/`, or `docs/COMMANDS.md` to learn usage.

**Do not** run `cards study` — study sessions are interactive and must be run by the user in the terminal.

| File | Purpose |
| --- | --- |
| [reference.md](reference.md) | Flags, JSON shapes, config keys |
| [examples.md](examples.md) | User phrase → command mapping |

## Setup

1. `cards` on PATH (`go install ./cmd/cards` from repo root)
2. Else `./cards` after `go build -o cards ./cmd/cards`

Always append `--json` for management commands.

## Commands

| Command | Notes |
| --- | --- |
| `cards deck create "<name>" --json` | |
| `cards deck list --json` | |
| `cards deck delete "<name>" --json --yes` | `--yes` required with `--json` |
| `cards add "<deck>" --front "..." --back "..." --json` | New cards go to front of queue |
| `cards list "<deck>" --json` | Metadata only, not queue order |
| `cards search "<query>" [--term "..."] [--deck "<deck>"] --json` | OR-matched terms across front, back, deck name |
| `cards show "<deck>" <id> --json` | |
| `cards edit "<deck>" <id> [--front "..."] [--back "..."] --json` | At least one of `--front` / `--back` |
| `cards delete "<deck>" <id> --json` | |
| `cards queue "<deck>" --json` | Queue order inspection |
| `cards stats "<deck>" --json` | Deck session count, last session, nudge |
| `cards export "<deck>" --format json --json` | Export summary (use `--output` for file) |
| `cards import --deck "<deck>" --format json --file path --json` | Import cards |
| `cards config --json` | Resolved paths and study defaults |
| `cards stats <deck>` | Available |
| `cards version --json` | |

## Errors

Exit `0` = success; `1` = validation, not found, or DB error. If `cards` missing: `go install ./cmd/cards` or build `./cards`.

See [reference.md](reference.md) for JSON shapes and [examples.md](examples.md) for phrase mappings.
