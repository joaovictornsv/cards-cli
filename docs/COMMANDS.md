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

## `cards queue <deck>`

Show current queue order (position, card id, front preview).

```bash
cards queue portuguese --json
```

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

## `cards version`

Show CLI version and build metadata.

```bash
cards version
cards version --json
cards --version
```

## See also

- [PROJECT_DRAFT.md](PROJECT_DRAFT.md) — full product spec and scheduling algorithm
- [NEXT_STEPS.md](../NEXT_STEPS.md) — post-v1 features (including deferred stats)
