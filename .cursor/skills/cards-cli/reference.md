# cards-cli reference

**Binary:** `cards` on PATH · `go install ./cmd/cards` · Dev: `go build -o cards ./cmd/cards`

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | JSON output (use for all agent operations) |
| `--help` | Command help |
| `--version` | CLI version |

## Configuration

DB path (first match): `CARDS_DB` env → `database` in `~/.config/cards/config.toml` → `~/.local/share/cards/cards.db`

Global config keys in `config.toml`:

| Key | Default | Description |
| --- | --- | --- |
| `batch_size` | `4` | Default study session batch size |
| `again_offset` | `2` | Queue re-insert offset for `again` grade |
| `nudge_threshold_days` | `3` | Days without study before `cards stats` shows a nudge |

## Study (user-run only)

Agents do not run `cards study`. Grades affect queue position: `again` → front + `again_offset`, `easy` → end, `replace` → end + sets `replace_eligible`. If a user shares study output, session log JSON may include `deck`, `batch_size`, `reviews`, and `status` (`complete` or `quit`).

## Commands (management)

| Command | Status |
| --- | --- |
| `deck create <name>` | Available |
| `deck list` | Available |
| `deck delete <name>` | Available (`--yes` required with `--json`) |
| `add <deck> --front --back` | Available |
| `list <deck>` | Available (`--replace-eligible` filter) |
| `search [query]` | Available (`--term`, `--deck`) |
| `show <deck> <id>` | Available |
| `edit <deck> <id>` | Available (`--replace-eligible` to clear flag) |
| `delete <deck> <id>` | Available |
| `queue <deck>` | Available |
| `stats <deck>` | Available |
| `export <deck>` | Available (`--format`, `--output`) |
| `import` | Available (`--deck`, `--format`, `--file`, `--append`) |
| `study <deck>` | Available (user-run only; agents must not invoke) |
| `config` | Available |
| `version` | Available |

## JSON shapes

**Config** (`cards config --json`):

```json
{
  "database_path": "/home/user/.local/share/cards/cards.db",
  "config_path": "/home/user/.config/cards/config.toml",
  "config_exists": false,
  "source": "default",
  "batch_size": 4,
  "again_offset": 2,
  "nudge_threshold_days": 3
}
```

**Version** (`cards version --json`):

```json
{
  "version": "0.0.0-dev",
  "commit": "unknown",
  "go_version": "go1.25.0"
}
```

**Deck list** (`cards deck list --json`):

```json
{
  "decks": [
    { "id": 1, "name": "portuguese", "card_count": 42, "created_at": "2026-07-09T12:00:00Z" }
  ]
}
```

**Card add** (`cards add <deck> --front "..." --back "..." --json`):

```json
{
  "id": 1,
  "deck_id": 1,
  "front": "What is saudade?",
  "back": "A deep emotional state of longing.",
  "created_at": "2026-07-09T12:00:00Z",
  "updated_at": "2026-07-09T12:00:00Z",
  "replace_eligible": false
}
```

**Card list** (`cards list <deck> --json`):

```json
{
  "deck": "portuguese",
  "cards": [
    { "id": 1, "front": "What is saudade?", "created_at": "...", "updated_at": "...", "replace_eligible": false }
  ],
  "total": 1
}
```

Filter flagged cards: `cards list <deck> --replace-eligible --json` (same shape, only flagged cards).

**Card search** (`cards search [query] --json`):

```json
{
  "cards": [
    {
      "id": 1,
      "deck": "portuguese",
      "front": "What is saudade?",
      "back": "A deep emotional state of longing."
    }
  ],
  "total": 1
}
```

Flags: repeatable `--term` (OR-matched), optional positional query, optional `--deck` to scope to one deck.

**Queue** (`cards queue <deck> --json`):

```json
{
  "deck": "portuguese",
  "queue": [
    { "position": 0, "id": 3, "front_preview": "What is saudade?" }
  ]
}
```

**Deck stats** (`cards stats <deck> --json`):

```json
{
  "deck": "portuguese",
  "sessions_count": 5,
  "last_session_at": "2026-07-18T12:00:00Z",
  "last_session_ago": "3 days ago",
  "nudge": "last session: 3 days ago — ready for a quick review?"
}
```

`last_session_at` is `null` when never studied. `nudge` is empty when no nudge applies.

**Study session log** (user-run; JSON printed after interactive output):

```json
{
  "deck": "portuguese",
  "batch_size": 4,
  "deck_size": 42,
  "status": "complete",
  "reviews": [
    { "card_id": 1, "front": "What is saudade?", "grade": "easy", "position": 1 }
  ]
}
```

`grade` may also be `"replace"` (same queue effect as `easy`, sets `replace_eligible` on the card).

**Export summary** (`cards export <deck> --json` with `--output`):

```json
{
  "deck": "portuguese",
  "format": "json",
  "card_count": 42,
  "output": "/path/to/portuguese.json"
}
```

**Import result** (`cards import --deck <name> --file <path> --json`):

```json
{
  "deck": "portuguese",
  "cards_imported": 10,
  "errors": ["row 3: card back is required"]
}
```

## Flagged cards workflow

1. User flags cards during study with `r` (replace grade), or you inspect existing flags: `cards list <deck> --replace-eligible --json`
2. Rewrite content: `cards edit <deck> <id> --front "..." --back "..." --json`
3. Clear flag after refresh: `cards edit <deck> <id> --replace-eligible=false --json`
