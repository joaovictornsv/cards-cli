# cards-cli

A personal CLI for flashcard decks and terminal study sessions. Queue-position scheduling (not date-based), SQLite storage, single binary, no server.

```bash
cards deck create portuguese
cards add portuguese --front "What is saudade?" --back "A deep emotional state of longing."
cards study portuguese   # interactive — run yourself in the terminal
```

## Commands

| Command | Description | Status |
| --- | --- | --- |
| `deck create` | Create a new deck | Available |
| `deck list` | List decks with card counts | Available |
| `deck delete` | Delete a deck and its cards | Available |
| `add` | Add a card to a deck | Available |
| `list` | List cards in a deck (metadata) | Available |
| `show` | Show one card | Available |
| `edit` | Edit card front/back | Available |
| `delete` | Remove a card from deck and queue | Available |
| `queue` | Show current queue order | Available |
| `study` | Interactive study session (`--limit`, `--json`) | Available — **user-run only** |
| `config` | Show resolved configuration | Available |
| `version` | Show CLI version and build metadata | Available |

Use `--json` on management commands for scripting. Full flag reference: [docs/COMMANDS.md](docs/COMMANDS.md). For AI agents: [`.cursor/skills/cards-cli/SKILL.md`](.cursor/skills/cards-cli/SKILL.md).

## Setup

**Requirements:** Go 1.25+

```bash
git clone https://github.com/joaovictornsv/cards-cli.git
cd cards-cli
go build -o cards ./cmd/cards
```

Pre-built binaries for linux/amd64 will be available on [GitHub Releases](https://github.com/joaovictornsv/cards-cli/releases) after the first release.

### Database path

1. `CARDS_DB` environment variable
2. `database` in `~/.config/cards/config.toml`
3. Default: `~/.local/share/cards/cards.db`

```toml
database = "/home/user/cards.db"
batch_size = 4
again_offset = 2
```

Run `cards config` to see which path is in use.

## Development

```bash
go test ./...
go build -o cards ./cmd/cards
```

Optional local pre-commit hooks (requires [lefthook](https://github.com/evilmartians/lefthook)):

```bash
lefthook install
```

Changes are tracked in [CHANGELOG.md](CHANGELOG.md). Post-v1 ideas: [NEXT_STEPS.md](NEXT_STEPS.md).

## License

MIT — see [LICENSE](LICENSE).
