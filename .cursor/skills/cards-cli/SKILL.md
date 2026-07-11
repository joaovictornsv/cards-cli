---
name: cards-cli
description: >-
  Manage flashcard decks and cards via the cards CLI (deck create/list/delete,
  add, list, show, edit, delete, queue, stats, config). Use when the user
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

## Deck management

```bash
cards deck create "<name>" --json
cards deck list --json
cards deck delete "<name>" --json
```

## Card management

```bash
cards add "<deck>" --front "..." --back "..." --json
cards list "<deck>" --json
cards show "<deck>" <id> --json
cards edit "<deck>" <id> [--front "..."] [--back "..."] --json
cards delete "<deck>" <id> --json
```

## Queue and stats

```bash
cards queue "<deck>" --json
cards stats "<deck>" --json
```

## Configuration

```bash
cards config --json
```

## Errors

Exit `0` = success; `1` = validation, not found, or DB error. If `cards` missing: `go install ./cmd/cards` or build `./cards`.

See [reference.md](reference.md) for JSON shapes and [examples.md](examples.md) for phrase mappings.

**Note:** Many commands are not yet implemented in early versions — check CLI help or repo issues if a command fails.
