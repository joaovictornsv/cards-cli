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

## Grades (study — user-run only)

| Grade | Queue behavior |
| --- | --- |
| `again` | Insert at front + `again_offset` |
| `hard` | Insert at front + `hard_offset` |
| `easy` | Insert at end of queue |

Agents do not run study sessions. Documented here for context when inspecting queue state.

## Commands (management)

| Command | Status |
| --- | --- |
| `deck create <name>` | Available |
| `deck list` | Available |
| `deck delete <name>` | Available |
| `add <deck> --front --back` | Available |
| `list <deck>` | Available |
| `show <deck> <id>` | Available |
| `edit <deck> <id>` | Available |
| `delete <deck> <id>` | Available |
| `queue <deck>` | Planned (v1) |
| `stats <deck>` | Planned (v1) |
| `config` | Available |
| `version` | Available |

## JSON shapes (stubs)

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

**Queue** (planned):

```json
{
  "deck": "portuguese",
  "queue": [
    { "position": 0, "id": 3, "front_preview": "What is saudade?" }
  ]
}
```

**Stats** (planned):

```json
{
  "deck": "portuguese",
  "card_count": 42,
  "sessions_run": 5,
  "last_session_at": "2026-07-08T18:30:00Z",
  "nudge": "last session: 1 day ago"
}
```
