# Playoff Bracket Fixes Design

**Date:** 2026-03-15

## Overview

Three fixes to the Playoff Bracket TUI view:
1. East/West sides swapped (East should be on the right)
2. Past seasons show all-TBD due to silent fallback to standings projection
3. Bracket ASCII renderer replaced with a cleaner columnar table view

---

## Issue 1: East/West Conference Sides

### Problem
The bracket renders East on the left and West on the right. NBA convention is West on the left, East on the right.

### Fix
Pure positional changes — no logic changes to series data or converter indices:

- `tui/playoffBracket.go`: update `colSeriesCount` comment and swap column assignments in `cursorIndexForColRow` and `colRowForIndex` so West occupies cols 0-2 and East occupies cols 4-6
- `tui/bracketTableRenderer.go` (Issue 3): assemble the rendered output as `west_rendered + finalsCol + east_rendered` instead of `east + finalsCol + west`. Data routing to each builder is unchanged: `series[0:7]` (East data) feeds the East-half builder, `series[8:15]` (West data) feeds the West-half builder. Only the final join order flips.

Note: the underlying series index order (0-3=East R1, 4-5=East Semis, 6=East Finals, 7=Finals, 8=West Finals, 9-10=West Semis, 11-14=West R1) is **not changed**. The `cursor == N` checks inside each half-builder continue to use the same absolute series indices — no changes needed there.

---

## Issue 2: Past Seasons Showing All-TBD

### Problem
`fetchPlayoffBracketCmd` silently swallows errors at three points in the fetch chain:

```go
if err := cl.FetchCommonPlayoffSeries(season); err == nil {
    if rs, err := cl.Loader.LoadCommonPlayoffSeries(season); err == nil {
        if bracket, err := converters.PopulatePlayoffBracket(rs, season); err == nil {
            return bracketFetchedMsg{bracket: bracket}
        }
    }
}
// Falls through to standings projection → all TBD for past seasons
```

The most likely root causes (in order of probability):
1. **`PopulatePlayoffBracket` returns an error** — it returns `"no playoff series data for season %s"` when `RowSet` is empty, which can happen if the API response format differs for historical seasons
2. **`FetchCommonPlayoffSeries` fails** — wrong request parameters or URL for historical season strings
3. **`LoadCommonPlayoffSeries` fails** — cache file not written or path mismatch

### Fix
1. Add `log.Printf` at each of the three inner failure points with season and error message
2. Run against a known completed season (e.g. 2023-24) and read `~/.config/nba-tui/logs/appLog.log`
3. Fix the identified root cause

Acceptance criteria: navigating to a completed past season (e.g. 2023-24) shows real team names and win counts, not TBD.

---

## Issue 3: Columnar Table Renderer

### Problem
The current ASCII bracket renderer (`bracketRenderer.go`) produces asymmetric, hard-to-read output that clips on narrow terminals.

### Solution
Replace with a columnar table renderer (`bracketTableRenderer.go`).

### Layout

```
     WEST                         │        │           EAST
 R1        Semis      CF          │ FINALS │  CF       Semis      R1
─────────────────────────────────────────────────────────────────────
(1)OKC 4                          │        │                      (1)BOS 4
(8)MEM 1  (1)OKC 4                │        │  (1)BOS 4  (1)BOS 4  (8)MIA 1
(4)LAL 2             (1)OKC 4     │        │            (4)NYK 2   (4)NYK 2
(5)DEN 1  (4)LAL 1                │OKC 4-2 │                       (5)ORL 1
(2)GSW 4                          │  BOS   │                       (2)BKN 4
(7)POR 0  (2)GSW 4   (2)GSW 2     │        │  (2)BKN 2  (2)BKN 4  (7)TOR 1
(3)HOU 3                          │        │            (3)MIL 3   (3)MIL 3
(6)MIN 2  (3)HOU 2                │        │                       (6)CHI 2
```

### Row Placement Rules (identical for both halves)

8 display lines total (2 per R1 slot, 4 slots). Finals is centered at lines 3-4.

| Line | W R1        | W Semis         | W CF      | Finals     | E CF      | E Semis         | E R1        |
|------|-------------|-----------------|-----------|------------|-----------|-----------------|-------------|
| 0    | R1[0] top   |                 |           |            |           |                 | R1[0] top   |
| 1    | R1[0] bot   | Semi[0] top     |           |            |           | Semi[0] top     | R1[0] bot   |
| 2    | R1[1] top   |                 | CF top    |            | CF top    |                 | R1[1] top   |
| 3    | R1[1] bot   | Semi[0] bot     |           | Finals top |           | Semi[0] bot     | R1[1] bot   |
| 4    | R1[2] top   |                 |           | Finals bot |           |                 | R1[2] top   |
| 5    | R1[2] bot   | Semi[1] top     | CF bot    |            | CF bot    | Semi[1] top     | R1[2] bot   |
| 6    | R1[3] top   |                 |           |            |           |                 | R1[3] top   |
| 7    | R1[3] bot   | Semi[1] bot     |           |            |           | Semi[1] bot     | R1[3] bot   |

### Cell Format
- Active series: `(seed)TRI wins` e.g. `(1)BOS 4`
- TBD series: `TBD` (dimmed style)
- Selected cell: highlighted style (both lines of the series)
- Empty cell: blank padded to column width

### Navigation
- Existing cursor model (`cursorCol`, `cursorRow`, `colSeriesCount`) unchanged
- `newBracketRenderer` call in `playoffBracket.go` renamed to `newBracketTableRenderer`

### Files Changed
- `tui/bracketRenderer.go` — deleted
- `tui/bracketRenderer_test.go` — deleted
- `tui/bracketTableRenderer.go` — new file (replaces bracketRenderer.go)
- `tui/bracketTableRenderer_test.go` — new file (replaces bracketRenderer_test.go); must cover: correct row count (8 lines), cell content at known positions (e.g. R1[0] top at line 0, CF top at line 2, Finals top at line 3), cursor highlighting applied to both lines of a series, and TBD cells rendered dimmed
- `tui/playoffBracket.go` — cursor column order updated (Issue 1) + logging added (Issue 2) + call site renamed

Note: `cmd/converters/playoffConverters.go` and `cmd/converters/playoffConverters_test.go` have pre-existing in-progress modifications unrelated to this spec. No converter changes are required by this spec.

---

## Out of Scope
- Wins data for active (in-progress) series — existing behaviour preserved
- Navigation key changes
- Changes to `PlayoffSeries`, `PlayoffBracket`, or converter types
