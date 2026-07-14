# cards-cli — Project Draft

> CLI flashcard app for terminal study sessions.  
> Scheduling model: **queue-position based** (not date/time based like Anki).  
> Status: **decisions settled** — ready for implementation.

---

## 1. Vision

A personal CLI tool to create decks, add cards, and run short study sessions in the terminal. Inspired by Anki’s review loop, but with a deliberately simpler scheduling strategy:

- The deck is an **ordered queue**.
- Reviewing a card **moves it** forward in the queue based on how hard it was.
- Each study session processes a **fixed batch** (default: 4 cards).
- No due dates, no intervals in days/hours — only queue position.

**Not a goal (v1):** replicate Anki, sync with Anki, rich media, mobile, or scientifically optimal long-term spaced repetition.

**Is a goal (v1):** fast daily drilling from the terminal, small sessions, simple mental model, single binary, local-first storage.

### 1.1 Usage model (human vs AI)

Unlike `books-cli` (primarily driven by AI agents — the user never runs those commands directly), `cards-cli` splits responsibilities:

| Who | What |
|---|---|
| **User (terminal)** | Run **interactive study sessions** (`cards study <deck>`) — one card at a time, grade with keyboard |
| **AI agents** | **Deck and card management** — create/update/delete decks and cards, inspect queue, etc. |

Study sessions must be **fully interactive** in the terminal: show one card at a time, let the user reveal the back, then choose a grade (`again`, `hard`, `easy`) via arrow keys and/or number keys.

---

## 2. Scheduling model (core algorithm)

### 2.1 Queue semantics

```
[ front = see sooner ]  card₃  card₇  card₁  ...  cardₙ  [ back = see later ]
```

- The deck is a single ordered list of card IDs.
- **Front** = highest priority (reviewed next).
- **Back** = lowest priority (reviewed last, after everything else cycles).

### 2.2 Study session

1. Take up to **N cards** from the **front** of the queue (default `N = 4`).
2. For each card: show front → user reveals back → user grades.
3. Re-insert the card into the queue **immediately** after it is graded (live queue updates).
4. Persist the new queue order after each graded card.
5. Session ends when the batch is done (or user quits early).

Sessions are independent of calendar day. Multiple sessions per day are allowed — including **two sessions in the same day for different decks**. Each session is **one deck only**.

### 2.3 Grading → re-insert rules (decided)

| Grade | Behavior | Intent |
|---|---|---|
| `again` | insert at position `front + 2` | see again very soon |
| `hard` | insert at position `front + 5` | see again relatively soon (fewer positions than `again` would mean sooner — confirmed) |
| `easy` | insert at **end** of queue | deprioritize for a long while |

**Decided:** “position” is always measured from the **front of the queue at re-insert time** (after the card was removed for review). Simpler and predictable.

**Decided:** Offsets are **global** (config), not per-deck in v1.

**Decided:** Same-session repeats are **allowed** — if `again` → `+2` and batch size is 4, the card may reappear within the same session.

### 2.4 What this model does *not* do

- No time-based forgetting (a card at the back does not “become due” after X days).
- Spacing in real-world time **emerges** from: deck size, batch size, and how often you run sessions.
- Skipping days does not change the queue — you resume where you left off.

---

## 3. CLI surface

Binary name: **`cards`**.

### 3.1 Commands (v1)

| Command | Description | Primary user |
|---|---|---|
| `cards deck create <name>` | Create a new deck | AI |
| `cards deck list` | List decks with card counts | AI |
| `cards deck delete <name>` | Delete a deck and its cards | AI |
| `cards add <deck> --front "..." --back "..."` | Add a card | AI |
| `cards list <deck>` | List cards (metadata only, not full queue walk) | AI |
| `cards show <deck> <id>` | Show one card | AI |
| `cards edit <deck> <id>` | Edit front/back | AI |
| `cards delete <deck> <id>` | Remove a card from deck and queue | AI |
| `cards queue <deck>` | Show current queue order (debug/inspection) | AI |
| `cards study <deck>` | Run an **interactive** study session | **User** |
| `cards config` | Show resolved config paths | AI |
| `cards version` | Version and build info | either |

No `cards study --all` or multi-deck session in v1.

### 3.2 `study` flags

```
cards study <deck> [--limit 4] [--json]
```

- `--limit` — batch size (default 4).
- `--json` — machine-readable session log (for scripting; secondary to interactive UX).

### 3.3 Study session UX (decided)

Interactive terminal session — one card at a time:

```
$ cards study portuguese

Session: portuguese (batch 4/4, 42 cards in deck)

[1/4] What is "saudade"?
      (↑↓ or 1/2/3 to grade after reveal, space/enter to reveal, q to quit)

→ reveal (space/enter)

      A Portuguese word for a deep emotional state of longing.

      [1] again   [2] hard   [3] easy
      (arrow keys or number keys)

... next card ...
```

**Decided:**

- **Interactive** — not a batch dump; user drives each card.
- Show **front** first; user **reveals back** (space/enter); back may also auto-show as part of the flow before grading.
- Grade with **arrow keys** and/or **number keys** (`1` = again, `2` = hard, `3` = easy).
- Progress indicator per batch (`[1/4]`, etc.).
- Plain text rendering only in v1 (no markdown).
- English CLI messages.

**Mid-session quit (`q`):** save progress after each graded card — reviewed cards stay re-inserted; unreviewed cards from the batch remain at front untouched.

---

## 4. Data model

### 4.1 Entities

**Deck**
- `id`, `name`, `created_at`
- optional: `description`

**Card**
- `id`, `deck_id`
- `front` (plain text)
- `back` (plain text)
- `created_at`, `updated_at`

**Queue**
- ordered list of `card_id` per deck
- stored as explicit `position` column or serialized order

### 4.2 Storage (decided)

Align with `books-cli` conventions:

- **SQLite** at `~/.local/share/cards/cards.db`
- Config at `~/.config/cards/config.toml`
- Override via `CARDS_DB` env var
- Single Go binary, no server

Export/import JSON deferred to v1.1.

**No relationship to `books-cli`** — independent config paths and data; integration is a future nice-to-have only.

### 4.3 New card insertion (decided)

New cards are inserted at the **front** of the queue so they are introduced immediately (not buried behind a large backlog).

### 4.4 Deletion (decided)

Deleting a card removes it from the deck and queue automatically. **No archive feature** in v1.

### 4.5 Duplicates (decided)

**Duplicate front text is allowed** in a deck. No mandatory `check` command for duplicates in v1.

---

## 5. Tech stack

Mirror `books-cli` since it’s a proven pattern in this environment:

| Layer | Choice |
|---|---|
| Language | Go 1.25+ |
| CLI framework | cobra |
| Storage | SQLite (modernc.org/sqlite or mattn/go-sqlite3) |
| Study UX | Interactive terminal (arrow/number keys); plain stdin/terminal control for v1 |
| Output | table + `--json` for scripting (management commands) |
| Repo | `github.com/joaovictornsv/cards-cli` (?) |

---

## 6. MVP scope

### In scope (v1)

- [ ] One or more decks
- [ ] Plain text front/back cards
- [ ] Queue-based scheduling with 3 grades (`again`, `hard`, `easy`)
- [ ] **Interactive** study sessions with configurable batch size (default 4)
- [ ] Persistent queue across sessions; immediate re-insert per graded card
- [ ] New cards enter at **front** of queue
- [ ] `cards config`, `cards version`, `--json` on list commands
- [ ] Mid-session quit saves per-card progress
- [ ] Empty deck → friendly error; small deck → `min(limit, deck_size)` cards

### Out of scope (v1)

- Anki import/export
- Images, audio, LaTeX
- Date-based scheduling / SM-2 / FSRS
- Sync, multi-device
- Cloze deletions
- Tags
- Web UI
- Archive for deleted/hidden cards
- Shuffle within batch
- Stale bump
- Multi-deck study session
- Markdown rendering
- Relationship / shared config with `books-cli`

### Nice to have (v1.1+)

- CSV import/export
- Tags and filtered study (`study --tag verbs`)
- **Shuffle** study order within a batch
- **Stale bump** — cards not seen in X sessions move toward front (still no dates)
- Integration hook with `books-cli` (generate cards from book notes)
- Markdown rendering for card content
- JSON export/import

---

## 7. Decisions log (resolved)

| # | Topic | Decision |
|---|---|---|
| Q1 | Project name / binary | **`cards`** |
| Q2 | Grading scale | **3 grades:** `again`, `hard`, `easy` |
| Q3 | Queue offsets | `again` = +2, `hard` = +5, `easy` = end; relative to **front at re-insert time**; global config |
| Q4 | Same-session repeats | **Allow** |
| Q5 | New cards in queue | **Front** of queue |
| Q6 | Empty / small deck | Study `min(limit, deck_size)`; empty deck → friendly error + hint |
| Q7 | Storage | **SQLite** |
| Q8 | Card content | **Plain text** |
| Q9 | Multiple decks per session | **One deck per session**; multiple sessions per day OK (even different decks) |
| Q10 | Default batch size | Global default in `config.toml` (4), overridable via `--limit` |
| Q11 | Quit mid-session | **Option A** — save after each graded card; unreviewed batch cards stay at front |
| Q12 | Stale bump | **Defer** to v1.1 |

---

## 8. Clarifications log (resolved)

| # | Topic | Clarification |
|---|---|---|
| U1 | Daily habit / nudges | **Deferred** — see [NEXT_STEPS.md](../NEXT_STEPS.md) for stats/nudge command |
| U2 | Original “+5 / +2” wording | **Confirmed:** fewer positions forward = sooner review (`again` +2 sooner than `hard` +5) |
| U3 | Archive on delete | **No** archive feature |
| U4 | Duplicate cards | **Allow** duplicate front text |
| U5 | Shuffle within batch | **Defer** to next version |
| U6 | Language | **English** CLI messages |
| U7 | Target deck size | **Small** (~50 cards) — no special pagination/perf work needed for v1 |
| U8 | books-cli relationship | **None** for v1 — fully independent |
| U9 | Re-insert order in batch | **Immediate re-insert** per card after grading (not batch re-insert at end) |

---

## 9. Example walkthrough (confirmed)

Deck queue (front → back): `[A, B, C, D, E, F, G, H]` (8 cards)  
Session limit: 4

**Step 1 — pull batch:** review A, B, C, D (removed from queue)  
Remaining queue: `[E, F, G, H]`

**Step 2 — grade and re-insert (immediate, per card):**

| Card | Grade | Insert at | Queue after insert |
|---|---|---|---|
| A | easy | end | `[E, F, G, H, A]` |
| B | again (+2) | index 2 | `[E, F, B, G, H, A]` |
| C | hard (+5) | index 5 | `[E, F, B, G, H, C, A]` |
| D | easy | end | `[E, F, B, G, H, C, A, D]` |

**Next session** starts with: `E, F, B, G`.

Matches intent.

---

## 10. Decision checklist

- [x] **Name:** `cards`
- [x] **Grades:** 3 (`again`, `hard`, `easy`)
- [x] **Offsets:** again=+2, hard=+5, easy=end (relative to front at re-insert)
- [x] **New cards:** front
- [x] **Same-session repeats:** allow
- [x] **Storage:** SQLite
- [x] **Card format:** plain text
- [x] **Sessions:** one deck per session
- [x] **Default batch size:** 4 (global config, `--limit` override)
- [x] **Mid-session quit:** save per card
- [x] **Stale bump:** defer
- [x] **Nudges:** deferred (see NEXT_STEPS.md)
- [x] **Study UX:** interactive, arrows/numbers, user-run
- [x] **Management:** AI-driven
- [x] **Shuffle batch:** defer
- [x] **Archive:** no
- [x] **books-cli tie-in:** none (v1)

---

## 11. Next steps

1. ~~Resolve all items in §7 and §8.~~ Done.
2. Create repo (`cards-cli`) and scaffold Go project (cobra + SQLite).
3. Implement data layer: decks, cards, queue persistence.
4. Implement interactive `study` with batch pull + reveal + grade + immediate re-insert.
5. Add management commands (`add`, `list`, `deck`, etc.) — optimized for AI agent use + `--json`.
6. Manual dogfooding with one real deck (~20 cards).
7. Write `docs/COMMANDS.md` and agent skill for **deck/card management** (study is user-run; skill documents management surface).

---

## 12. Revision log

| Date | Change |
|---|---|
| 2026-07-09 | Initial draft from conversation |
| 2026-07-09 | All §7/§8 decisions recorded; usage model (AI management / human study); interactive UX spec |
