# Dynamic Table Height Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix table header clipping by dynamically sizing tables based on terminal height.

**Architecture:** Add a helper function to calculate page size from terminal dimensions, then apply it in each view's Update() handler when WindowSizeMsg arrives.

**Tech Stack:** Go, bubble-table (WithPageSize API), bubbletea (WindowSizeMsg)

---

### Task 1: Add calculatePageSize Helper

**Files:**
- Modify: `tui/stylesAndMappings.go` (add function at end)
- Create: `tui/stylesAndMappings_test.go`

**Step 1: Write the failing test**

Create `tui/stylesAndMappings_test.go`:

```go
package tui

import "testing"

func TestCalculatePageSize_SingleTable(t *testing.T) {
	// Terminal height 40, single table
	// Available = 40 - 4(margins) - 2(help) - 1(header) - 2(border) - 1(buffer) = 30
	result := calculatePageSize(40, 1)
	if result != 30 {
		t.Errorf("expected 30, got %d", result)
	}
}

func TestCalculatePageSize_DualTable(t *testing.T) {
	// Terminal height 40, dual tables
	// Available = 40 - 4(margins) - 2(help) - 2(labels) - 2(buffer) = 30
	// Per table = (30 - 6(overhead per table)) / 2 = 12
	result := calculatePageSize(40, 2)
	if result != 12 {
		t.Errorf("expected 12, got %d", result)
	}
}

func TestCalculatePageSize_MinimumEnforced(t *testing.T) {
	// Tiny terminal should return minimum of 3
	result := calculatePageSize(10, 2)
	if result != 3 {
		t.Errorf("expected minimum 3, got %d", result)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./tui -run TestCalculatePageSize -v`
Expected: FAIL with "undefined: calculatePageSize"

**Step 3: Write minimal implementation**

Add to end of `tui/stylesAndMappings.go`:

```go
// calculatePageSize returns the number of rows a table should display
// based on terminal height and number of tables in the view.
func calculatePageSize(terminalHeight, tableCount int) int {
	const (
		margins    = 4 // DocStyle margin top + bottom
		helpFooter = 2
		buffer     = 1
	)

	var pageSize int

	if tableCount == 1 {
		// Single table: header(1) + border(2)
		overhead := margins + helpFooter + 1 + 2 + buffer
		pageSize = terminalHeight - overhead
	} else {
		// Dual tables: labels(2) + extra buffer(1)
		baseOverhead := margins + helpFooter + 2 + 2
		available := terminalHeight - baseOverhead
		// Each table needs header(1) + border(2) = 3
		perTableOverhead := 3
		pageSize = (available - (perTableOverhead * tableCount)) / tableCount
	}

	if pageSize < 3 {
		pageSize = 3
	}
	return pageSize
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./tui -run TestCalculatePageSize -v`
Expected: PASS

**Step 5: Commit**

```bash
git add tui/stylesAndMappings.go tui/stylesAndMappings_test.go
git commit -m "feat: add calculatePageSize helper for dynamic table sizing"
```

---

### Task 2: Apply Dynamic Page Size in LeagueLeaders

**Files:**
- Modify: `tui/leagueLeaders.go:99-109` (WindowSizeMsg handler)

**Step 1: Locate current WindowSizeMsg handler**

Current code at line 99-109:
```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}
```

**Step 2: Update handler to apply page size**

Replace the WindowSizeMsg case with:

```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}

	pageSize := calculatePageSize(msg.Height, 1)
	m.leaderboard = m.leaderboard.WithPageSize(pageSize)
```

**Step 3: Run existing tests**

Run: `go test ./tui -v`
Expected: PASS (no behavior change to existing tests)

**Step 4: Manual verification**

Run: `go build -o ./bin/nba-now-linux ./cmd/main.go && ./bin/nba-now-linux 2026-02-12`
Navigate to League Leaders, resize terminal - table should adjust row count.

**Step 5: Commit**

```bash
git add tui/leagueLeaders.go
git commit -m "feat: apply dynamic page size to leagueLeaders view"
```

---

### Task 3: Apply Dynamic Page Size in SeasonStandings

**Files:**
- Modify: `tui/seasonStandings.go:152-161` (WindowSizeMsg handler)

**Step 1: Locate current WindowSizeMsg handler**

Current code at line 152-161:
```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}
```

**Step 2: Update handler to apply page size**

Replace the WindowSizeMsg case with:

```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}

	pageSize := calculatePageSize(msg.Height, 2)
	m.eastTeams = m.eastTeams.WithPageSize(pageSize)
	m.westTeams = m.westTeams.WithPageSize(pageSize)
```

**Step 3: Run existing tests**

Run: `go test ./tui -v`
Expected: PASS

**Step 4: Manual verification**

Run: `./bin/nba-now-linux 2026-02-12`
Navigate to Season Standings - both tables should fit within terminal height.

**Step 5: Commit**

```bash
git add tui/seasonStandings.go
git commit -m "feat: apply dynamic page size to seasonStandings view"
```

---

### Task 4: Apply Dynamic Page Size in BoxScore

**Files:**
- Modify: `tui/boxScore.go:301-311` (WindowSizeMsg handler)

**Step 1: Locate current WindowSizeMsg handler**

Current code at line 301-311:
```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}
```

**Step 2: Update handler to apply page size**

Replace the WindowSizeMsg case with:

```go
case tea.WindowSizeMsg:
	m.width = msg.Width
	m.height = msg.Height

	if m.width > m.maxWidth {
		m.maxWidth = m.width
	}
	if m.height > m.maxHeight {
		m.maxHeight = m.height
	}

	pageSize := calculatePageSize(msg.Height, 2)
	m.homeTeamBoxScore = m.homeTeamBoxScore.WithPageSize(pageSize)
	m.awayTeamBoxScore = m.awayTeamBoxScore.WithPageSize(pageSize)
```

**Step 3: Run all tests**

Run: `go test ./... -v`
Expected: PASS

**Step 4: Manual verification**

Run: `./bin/nba-now-linux 2026-02-12`
Navigate to Daily View → select a game → Box Score should fit within terminal.

**Step 5: Commit**

```bash
git add tui/boxScore.go
git commit -m "feat: apply dynamic page size to boxScore view"
```

---

### Task 5: Final Verification & Cleanup

**Step 1: Run full test suite**

Run: `make test`
Expected: All tests PASS

**Step 2: Build all platforms**

Run: `make build`
Expected: Binaries created in ./bin/

**Step 3: End-to-end verification**

Test in Hyprland tiled mode:
1. Launch app with a recent game date
2. Navigate to Season Standings - both tables visible, headers not clipped
3. Navigate to a Box Score - both tables visible, headers not clipped
4. Resize terminal - tables should dynamically adjust

**Step 4: Final commit (if any cleanup needed)**

```bash
git add -A
git commit -m "chore: cleanup after dynamic table height implementation"
```
