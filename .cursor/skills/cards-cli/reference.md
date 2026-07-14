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
| `hard_offset` | `5` | Queue re-insert offset for `hard` grade |

## Study (user-run only)

Agents do not run `cards study`. Grades affect queue position: `again` → front + `again_offset`, `hard` → front + `hard_offset`, `easy` → end. If a user shares study output, session log JSON may include `deck`, `batch_size`, `reviews`, and `status` (`complete` or `quit`).

## Commands (management)

| Command | Status |
| --- | --- |
| `deck create <name>` | Available |
| `deck list` | Available |
| `deck delete <name>` | Available (`--yes` required with `--json`) |
| `add <deck> --front --back` | Available |
| `list <deck>` | Available |
| `show <deck> <id>` | Available |
| `edit <deck> <id>` | Available |
| `delete <deck> <id>` | Available |
| `queue <deck>` | Available |
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
  "hard_offset": 5
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
  "updated_at": "2026-07-09T12:00:00Z"
}
```

**Card list** (`cards list <deck> --json`):

```json
{
  "deck": "portuguese",
  "cards": [
    { "id": 1, "front": "What is saudade?", "created_at": "...", "updated_at": "..." }
  ],
  "total": 1
}
```

**Queue** (`cards queue <deck> --json`):

```json
{
  "deck": "portuguese",
  "queue": [
    { "position": 0, "id": 3, "front_preview": "What is saudade?" }
  ]
}
```

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
