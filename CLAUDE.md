# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
make              # Build for Linux, Windows, macOS (arm64)
make test         # Run all tests (go test ./...)
make clean        # Remove binaries

go test ./cmd/converters/...      # Run converter tests only
go test ./tui/...                 # Run TUI tests only
go test -v -run TestName ./...    # Run specific test
```

## Running the App

```bash
./bin/nba-now-linux YYYY-MM-DD    # Date argument required (game date in local time)
```

## Architecture

**TUI Layer** (`tui/`) - Bubble Tea models for each view. Views communicate via messages, not direct calls. Each view implements `Update()` and `View()` methods.

**API Layer** (`cmd/nba/nbaAPI/`) - Interface-driven design:
- `HTTPRequester` - HTTP GET with NBA-specific headers
- `RequestBuilder` - URL construction from parameter structs
- `Client` - Orchestrates API calls, filesystem ops, data loading

**Data Flow**: API response → JSON file cache → DataLoader → Converters → TUI model

**Converters** (`cmd/converters/`) - Transform raw API JSON to domain types using reflection and struct tags (`isVisible`, `percentage`, `display`, `width`).

**Filesystem** (`cmd/nba/filesystem/`) - Caches JSON to `~/.config/nba-tui/`. Auto-cleans files >48h old on startup.

## Key Patterns

- **Factory functions** for dependency injection: `NewClient()`, `PathFactory()`, `NewDataLoader()`
- **Concurrent API calls** via goroutines + channels in `MakeDefaultRequests()` and `FetchTeamProfile()`
- **Struct tags control rendering** - converters use reflection to read display metadata
- **Custom Bubble Tea messages** for async operation coordination (e.g., `requestsFinishedMsg`)

## External Dependencies

- `charmbracelet/bubbletea` - TUI framework (Elm architecture)
- `charmbracelet/bubbles` - UI components
- `charmbracelet/lipgloss` - Terminal styling
- `evertras/bubble-table` - Table widget
- NBA Stats API: `https://stats.nba.com/stats/` (uses Scoreboard V3)

## Config Paths

- Cache: `~/.config/nba-tui/` (boxscores/, teamprofiles/, teamplayers/, news/)
- Logs: `~/.config/nba-tui/logs/appLog.log`
