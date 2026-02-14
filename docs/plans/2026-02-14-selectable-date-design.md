# Selectable Date — Design Document

## Summary

Replace the mandatory `-d YYYY-MM-DD` CLI argument with a default date (today in US Eastern) and an interactive date selector in the Daily Scores view. Users can browse dates with arrow keys or type a specific date.

## Default Date

- Remove the `-d` flag and all CLI argument parsing
- `GetCurrentDate()` returns today's date in `America/New_York` timezone
- No arguments required to launch the application

## Date Selector Component

New `DateSelector` struct embedded in `DailyView`, displayed at the top:

```
          ◀  2026-02-14  ▶
```

### Focus Zones

`DailyView` has two focus zones:
1. **Date selector** — default focus on entering Daily Scores
2. **Game cards grid** — Tab or Down switches focus; Tab switches back

### Navigation Mode (default)

- **Left arrow**: previous day (-1)
- **Right arrow**: next day (+1), capped at today (US Eastern)
- **Enter**: switch to edit mode (activates text input)
- **Tab / Down**: move focus to game cards

### Edit Mode

- Uses `charmbracelet/bubbles/textinput` (existing dependency)
- User types a date in YYYY-MM-DD format
- **Enter**: validate → fetch (blocking). Invalid input shows inline error, stays in edit mode
- **Escape**: cancel, revert to previous date, return to navigation mode

### Validation

- `time.Parse("2006-01-02", input)` for format validation
- Reject future dates (after today in US Eastern)
- Invalid input: show error inline, stay in edit mode

## Data Flow on Date Change

1. Date change triggers a `tea.Cmd` calling `Client.FetchDailyScoresForDate(date string)`
2. New explicit method on `Client` — does not mutate `DateProvider` shared state
3. Reuses existing pipeline: fetch JSON → cache to `~/.config/nba-tui/{date}_dsb` → deserialize
4. Returns `dailyScoresFetchedMsg` (same message type already used)
5. `DailyView` rebuilds game cards from new data
6. Input is blocked during fetch (no spinner, just wait)

## Error Handling

- **Invalid date**: inline error below selector, auto-dismiss on next keypress
- **API/network failure**: show "No games happened during YYYY-MM-DD", log to appLog.log
- **Rapid arrow keys**: blocked by input-blocking fetch behavior

## Files Affected

| File | Change |
|------|--------|
| `cmd/nba/nbaAPI/requestParameters.go` | Replace CLI arg parsing with US Eastern default |
| `cmd/nba/nbaAPI/client.go` | Add `FetchDailyScoresForDate(date)` method |
| `tui/dailyView.go` | Embed `DateSelector`, dual focus zones, new key handling |
| `tui/dateSelector.go` | New file: `DateSelector` struct with nav/edit modes |
| `tui/stylesAndMappings.go` | Possibly add date-related styles |
| `cmd/main.go` | Remove `os.Args` passing (if applicable) |
| `tui/mainMenu.go` | Update title to use new default date |

## Decisions Made

- **Hybrid UX**: arrow keys for quick browsing + text input for arbitrary jumps
- **Approach 1 chosen**: single `textinput` component with dual modes over digit-by-digit or three-field alternatives
- **Explicit fetch method** (Option A) over mutable DateProvider (Option B)
- **Blocking fetch** over async loading state
- **No future dates** allowed
- **No CLI arguments** — always defaults to today (US Eastern)
