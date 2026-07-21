# cards-cli command reference

Binary: **`cards`**. Use `--json` on management commands for machine-readable output.

All commands below are implemented.

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | JSON output (recommended for AI agents) |
| `--help` | Command help |
| `--version` | Print version and exit |

## Deck commands

### `cards deck create <name>`

Create a new deck.

```bash
cards deck create portuguese --json
```

### `cards deck list`

List all decks with card counts.

```bash
cards deck list --json
```

### `cards deck delete <name>`

Delete a deck and all its cards and queue entries.

| Flag | Description |
| --- | --- |
| `--yes`, `-y` | Confirm deletion without prompting (required with `--json`) |

```bash
cards deck delete portuguese --json --yes
```

### `cards deck shuffle <name>`

Randomly permute all card positions in the deck queue. Same cards remain in the deck; only order changes. Decks with 0–1 cards are a no-op (`status: "noop"`).

| Flag | Description |
| --- | --- |
| `--yes`, `-y` | Confirm shuffle without prompting (required with `--json`) |
| `--seed` | Deterministic shuffle seed (hidden; for testing) |

```bash
cards deck shuffle portuguese --json --yes
```

## Card commands

### `cards add <deck> --front "..." --back "..."`

Add a card to a deck. New cards are inserted at the **front** of the queue.

| Flag | Required | Description |
| --- | --- | --- |
| `--front` | yes | Card front (plain text) |
| `--back` | yes | Card back (plain text) |

```bash
cards add portuguese --front "What is saudade?" --back "A deep emotional state of longing." --json
```

### `cards list <deck>`

List cards in a deck (metadata: id, front text, timestamps, `replace_eligible`). Does not walk the full queue order.

| Flag | Description |
| --- | --- |
| `--replace-eligible` | List only cards flagged for content replacement |

```bash
cards list portuguese --json
cards list portuguese --replace-eligible --json
```

### `cards search [query]`

Search cards across all decks (or one deck with `--deck`). Matches card front, card back, and deck name (case-insensitive substring). Multiple `--term` flags are OR-matched; a positional query counts as one term.

| Flag | Description |
| --- | --- |
| `--term` | Search term (repeatable; terms are OR-matched) |
| `--deck` | Limit search to one deck |

```bash
cards search "saudade" --json
cards search --term "hello" --term "saudade" --json
cards search "hello" --deck portuguese --json
```

### `cards show <deck> <id>`

Show one card (full front and back).

```bash
cards show portuguese 3 --json
```

### `cards edit <deck> <id>`

Edit a card's front and/or back, or clear the `replace_eligible` flag. At least one of `--front`, `--back`, or `--replace-eligible` is required. Editing front/back does **not** clear `replace_eligible`.

| Flag | Description |
| --- | --- |
| `--front` | New front text |
| `--back` | New back text |
| `--replace-eligible` | Set flag (`true` or `false`; use `--replace-eligible=false` to clear) |

```bash
cards edit portuguese 3 --front "Updated question" --json
cards edit portuguese 3 --replace-eligible=false --json
```

### `cards delete <deck> <id>`

Remove a card from the deck and queue. No archive.

```bash
cards delete portuguese 3 --json
```

### `cards export <deck>`

Export a deck and all its cards in queue order.

| Flag | Default | Description |
| --- | --- | --- |
| `--format` | `json` | Export format: `json` or `csv` |
| `--output`, `-o` | stdout | Write export payload to file instead of stdout |
| `--json` | off | Print summary (`deck`, `format`, `card_count`, `output`) instead of mixing payload with summary |

**JSON format:** `{ "deck": "...", "cards": [{ "front": "...", "back": "...", "id": 1 }] }` — `id` is optional metadata for round-trip reference.

**CSV format:** header `front,back`, one card per line. Deck name is the positional `<deck>` argument.

```bash
cards export portuguese --format json
cards export portuguese --format csv -o portuguese.csv
cards export portuguese --format json -o portuguese.json --json
```

### `cards import`

Import cards from a JSON or CSV file. Creates the deck if it does not exist. Without `--append`, fails if the deck already exists.

| Flag | Required | Default | Description |
| --- | --- | --- | --- |
| `--deck` | yes | | Target deck name |
| `--format` | no | `json` | Import format: `json` or `csv` |
| `--file`, `-f` | yes | | Input file path (`-` for stdin) |
| `--append` | no | off | Add cards to an existing deck |
| `--json` | no | off | Summary output: `deck`, `cards_imported`, `errors` |

Imported cards are inserted at the **front** of the queue (same as `cards add`). Duplicate fronts are allowed.

```bash
cards import --deck portuguese --format json --file portuguese.json --json
cards import --deck portuguese --format csv --file portuguese.csv --append --json
```

## `cards queue <deck>`

Show current queue order (position, card id, front preview).

```bash
cards queue portuguese --json
```

## `cards stats <deck>`

Show per-deck study activity: session count, last session time, and an optional nudge when the deck has not been studied recently.

```bash
cards stats portuguese
cards stats portuguese --json
```

**JSON fields:** `deck`, `sessions_count`, `last_session_at` (RFC3339 or `null`), `last_session_ago` (e.g. `never`, `today`, `3 days ago`), `nudge` (empty when no nudge).

**Nudge threshold:** `nudge_threshold_days` in config (default `3`). A nudge appears when the deck was never studied or the last session is older than the threshold.

**Session counting:** Stats increment when a study session ends with status `complete`, or `quit` after at least one card was graded. Quitting before grading any card does not count as a session.

## `cards study <deck>`

Run an **interactive** study session. One card at a time: show front → reveal back → grade (`again`, `easy`, or `replace`).

**Primary user:** human in the terminal (not AI agents).

| Flag | Default | Description |
| --- | --- | --- |
| `--limit` | from config (`batch_size`, default 4) | Batch size for this session |
| `--json` | off | Machine-readable session log (printed after interactive output) |

**Grading → queue re-insert:**

| Grade | Behavior |
| --- | --- |
| `again` | Insert at front + `again_offset` (default 2) |
| `easy` | Insert at end of queue |
| `replace` | Same as `easy`; also sets `replace_eligible = true` on the card |

**Interactive controls:**

- Space/Enter — reveal back
- Arrow keys or `1`/`2` — grade (again / easy)
- `r` / `R` — replace (flag card for content refresh)
- `q` — quit mid-session (graded cards saved; unreviewed batch cards stay at front)

**Session stats:** See `cards stats` for counting rules and nudge behavior.

```bash
cards study portuguese
cards study portuguese --limit 6
```

## `cards config`

Show resolved configuration paths and study defaults.

```bash
cards config
cards config --json
```

**Database path resolution order:**

1. `CARDS_DB` environment variable
2. `database` key in `~/.config/cards/config.toml`
3. Default: `~/.local/share/cards/cards.db`

**Config file keys:**

| Key | Default | Description |
| --- | --- | --- |
| `database` | (see above) | SQLite database path |
| `batch_size` | `4` | Default study batch size |
| `again_offset` | `2` | Queue offset for `again` grade |
| `nudge_threshold_days` | `3` | Days without study before `cards stats` shows a nudge |

## `cards version`

Show CLI version and build metadata.

```bash
cards version
cards version --json
cards --version
```

## See also

- [PROJECT_DRAFT.md](PROJECT_DRAFT.md) — historical product spec and scheduling algorithm
- [NEXT_STEPS.md](../NEXT_STEPS.md) — future ideas and out-of-scope items
