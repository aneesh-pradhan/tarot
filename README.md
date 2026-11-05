# Tarot

A terminal tarot reading game built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

- The full 78-card deck (22 Major Arcana, 56 Minor Arcana), each with upright and reversed meanings
- Five spreads: Single Card, Past · Present · Future, Situation · Action · Outcome, Mind · Body · Spirit, and the ten-card Celtic Cross
- Ask an optional question, watch the deck shuffle, then turn cards over one at a time
- Reversed cards can be switched on or off
- The layout adapts to your terminal size, down to 80×24

## Install

```sh
go install github.com/aneesh-pradhan/tarot@latest
```

Or from a clone:

```sh
go run .
```

## Usage

```sh
tarot                 # random reading
tarot -seed 42        # reproducible reading
tarot -no-reversals   # upright cards only
```

| Key | Action |
| --- | --- |
| `↑`/`↓` or `1`–`5` | choose a spread |
| `r` | toggle reversed cards (menu) |
| `space` / `enter` | turn over the next card |
| `←`/`→`/`↑`/`↓` | select a card to read its meaning |
| `a` | reveal all cards |
| `n` / `esc` | new reading |
| `q` / `ctrl+c` | quit |

## Development

```sh
go test ./...
```
