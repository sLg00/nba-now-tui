# Dynamic Table Height Design

## Problem

Table headers in box scores and standings views get clipped ("eaten") at the top of the terminal when running in constrained window modes (e.g., Hyprland tiled mode). The content height exceeds terminal height, pushing the top off-screen.

## Solution

Dynamically calculate table page size based on terminal height. Tables will display only as many rows as fit within the available space, with the rest accessible via scrolling.

## Height Calculation

### Constants (in terminal rows)

| Element | Rows |
|---------|------|
| DocStyle margin (top + bottom) | 4 |
| Help footer | 2 |
| Table header | 1 per table |
| Table border (top + bottom) | 2 per table |
| Section labels (standings only) | 2 |
| Buffer | 1-2 |

### Formulas

**Single-table views** (leagueLeaders):
```
availableRows = terminalHeight - 4 - 2 - 1 - 2 - 1
pageSize = availableRows
```

**Dual-table views** (boxScore, standings):
```
availableRows = terminalHeight - 4 - 2 - 2 - 2
perTableRows = (availableRows - 6) / 2
pageSize = perTableRows
```

### Edge Cases

- Minimum page size: 3 rows (header visibility + 2 data rows)
- Initial render default: 10 rows (before first WindowSizeMsg)

## Implementation

### New Helper Function

Add to `tui/stylesAndMappings.go`:

```go
func calculatePageSize(terminalHeight, tableCount int) int
```

### Application Points

Page size applied in `Update()` when handling `tea.WindowSizeMsg`:

```go
case tea.WindowSizeMsg:
    m.height = msg.Height
    pageSize := calculatePageSize(msg.Height, 2)
    m.table1 = m.table1.WithPageSize(pageSize)
    m.table2 = m.table2.WithPageSize(pageSize)
```

## Files Changed

| File | Change |
|------|--------|
| `tui/stylesAndMappings.go` | Add `calculatePageSize()` helper |
| `tui/boxScore.go` | Apply dynamic page size on `WindowSizeMsg` |
| `tui/seasonStandings.go` | Apply dynamic page size on `WindowSizeMsg` |
| `tui/leagueLeaders.go` | Replace hardcoded `WithPageSize(20)` |

## Not Changing

- Table data structures
- Navigation/keybindings
- View layout structure
- Other files
