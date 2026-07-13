# cards-cli command reference

Binary: **`cards`**. Use `--json` on management commands for machine-readable output.

> **Scaffold note:** `config`, `version`, and `deck` commands are implemented. Other commands below are the v1 specification — marked where not yet available.

## Global flags

| Flag | Description |
| --- | --- |
| `--json` | JSON output (recommended for AI agents) |
| `--help` | Command help |
| `--version` | Print version and exit |

---

## `cards deck create <name>`

Create a new deck.

**Status:** Implemented.

| Flag | Description |
| --- | --- |
| `--json` | JSON output |

**Example:**

```bash
cards deck create portuguese --json
```

---

## `cards deck list`

List all decks with card counts.

**Status:** Implemented.

**Example:**

```bash
cards deck list --json
```

---

## `cards deck delete <name>`

Delete a deck and all its cards and queue entries.

**Status:** Implemented.

| Flag | Description |
| --- | --- |
| `--yes`, `-y` | Confirm deletion without prompting (required with `--json`) |
| `--json` | JSON output |

**Example:**

```bash
cards deck delete portuguese --json --yes
```

---

## `cards add <deck> --front "..." --back "..."`

Add a card to a deck. New cards are inserted at the **front** of the queue.

**Status:** Implemented.

| Flag | Required | Description |
| --- | --- | --- |
| `--front` | yes | Card front (plain text) |
| `--back` | yes | Card back (plain text) |
| `--json` | no | JSON output |

**Example:**

```bash
cards add portuguese --front "What is saudade?" --back "A deep emotional state of longing." --json
```

---

## `cards list <deck>`

List cards in a deck (metadata: id, front text, timestamps). Does not walk the full queue order.

**Status:** Implemented.

**Example:**

```bash
cards list portuguese --json
```

---

## `cards show <deck> <id>`

Show one card (full front and back).

**Status:** Not yet implemented.

**Example:**

```bash
cards show portuguese 3 --json
```

---

## `cards edit <deck> <id>`

Edit a card's front and/or back.

**Status:** Not yet implemented.

| Flag | Description |
| --- | --- |
| `--front` | New front text |
| `--back` | New back text |
| `--json` | JSON output |

At least one of `--front` or `--back` is required.

**Example:**

```bash
cards edit portuguese 3 --front "Updated question" --json
```

---

## `cards delete <deck> <id>`

Remove a card from the deck and queue. No archive.

**Status:** Not yet implemented.

**Example:**

```bash
cards delete portuguese 3 --json
```

---

## `cards queue <deck>`

Show current queue order (position, card id, front preview). Useful for debugging and agent inspection.

**Status:** Not yet implemented.

**Example:**

```bash
cards queue portuguese --json
```

---

## `cards study <deck>`

Run an **interactive** study session. One card at a time: show front → reveal back → grade (`again`, `hard`, `easy`).

**Primary user:** human in the terminal (not AI agents).

**Status:** Not yet implemented.

| Flag | Default | Description |
| --- | --- | --- |
| `--limit` | from config (`batch_size`, default 4) | Batch size for this session |
| `--json` | off | Machine-readable session log (secondary to interactive UX) |

**Grading → queue re-insert:**

| Grade | Behavior |
| --- | --- |
| `again` | Insert at front + `again_offset` (default 2) |
| `hard` | Insert at front + `hard_offset` (default 5) |
| `easy` | Insert at end of queue |

**Interactive controls:**

- Space/Enter — reveal back
- Arrow keys or `1`/`2`/`3` — grade (again / hard / easy)
- `q` — quit mid-session (graded cards saved; unreviewed batch cards stay at front)

**Example:**

```bash
cards study portuguese
cards study portuguese --limit 6
```

---

## `cards stats <deck>`

Deck statistics: card count, sessions run, last session time, nudge message.

**Status:** Not yet implemented.

**Example:**

```bash
cards stats portuguese --json
```

---

## `cards config`

Show resolved configuration paths and study defaults.

**Status:** Implemented.

**Example:**

```bash
cards config
cards config --json
```

**Resolution order for database path:**

1. `CARDS_DB` environment variable
2. `database` key in `~/.config/cards/config.toml`
3. Default: `~/.local/share/cards/cards.db`

**Config file keys:**

| Key | Default | Description |
| --- | --- | --- |
| `database` | (see above) | SQLite database path |
| `batch_size` | `4` | Default study batch size |
| `again_offset` | `2` | Queue offset for `again` grade |
| `hard_offset` | `5` | Queue offset for `hard` grade |

---

## `cards version`

Show CLI version and build metadata.

**Status:** Implemented.

**Example:**

```bash
cards version
cards version --json
cards --version
```

---

## See also

- [PROJECT_DRAFT.md](PROJECT_DRAFT.md) — full product spec and scheduling algorithm
- [NEXT_STEPS.md](../NEXT_STEPS.md) — post-v1 features
