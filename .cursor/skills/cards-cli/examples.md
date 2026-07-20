# cards-cli examples

Phrase → command. Always append `--json` for management operations. See [SKILL.md](SKILL.md) for agent rules.

**Study sessions (`cards study`) are user-run only — not listed here.**

## Decks

| User says | Command |
| --- | --- |
| "Create a deck called portuguese" | `cards deck create "portuguese" --json` |
| "List my decks" | `cards deck list --json` |
| "Delete the portuguese deck" | `cards deck delete "portuguese" --json --yes` |

## Cards

| User says | Command |
| --- | --- |
| "Add a card to portuguese" | `cards add "portuguese" --front "..." --back "..." --json` |
| "List cards in portuguese" | `cards list "portuguese" --json` |
| "Search for saudade across decks" | `cards search "saudade" --json` |
| "Search hello in portuguese deck" | `cards search "hello" --deck "portuguese" --json` |
| "List cards flagged for replacement" | `cards list "portuguese" --replace-eligible --json` |
| "Show card 3 in portuguese" | `cards show "portuguese" 3 --json` |
| "Edit card 3 front text" | `cards edit "portuguese" 3 --front "new front" --json` |
| "Clear replace flag on card 3" | `cards edit "portuguese" 3 --replace-eligible=false --json` |
| "Delete card 3 from portuguese" | `cards delete "portuguese" 3 --json` |

## Queue

| User says | Command |
| --- | --- |
| "Show the queue for portuguese" | `cards queue "portuguese" --json` |

## Config

| User says | Command |
| --- | --- |
| "Where is my cards database?" | `cards config --json` |
| "What version of cards?" | `cards version --json` |
