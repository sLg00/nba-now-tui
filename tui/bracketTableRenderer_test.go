package tui

import (
    "strings"
    "testing"

    "github.com/sLg00/nba-now-tui/cmd/nba/types"
)

func makeTableSeries(tricode string, seed, wins int) types.PlayoffSeries {
    return types.PlayoffSeries{
        Status:     "active",
        TopTeam:    types.PlayoffTeam{Tricode: tricode + "T", Seed: seed, Wins: wins},
        BottomTeam: types.PlayoffTeam{Tricode: tricode + "B", Seed: seed + 7, Wins: 0},
    }
}

func makeFullBracket() []types.PlayoffSeries {
    s := make([]types.PlayoffSeries, 15)
    // East R1: indices 0-3
    s[0] = makeTableSeries("E1", 1, 3) // 1v8
    s[1] = makeTableSeries("E4", 4, 2) // 4v5
    s[2] = makeTableSeries("E3", 3, 1) // 3v6
    s[3] = makeTableSeries("E2", 2, 0) // 2v7
    // East Semis: indices 4-5
    s[4] = makeTableSeries("ES", 1, 2)
    s[5] = makeTableSeries("EL", 2, 1)
    // East Finals: index 6
    s[6] = makeTableSeries("EC", 1, 1)
    // Finals: index 7
    s[7] = makeTableSeries("FN", 1, 2)
    // West Finals: index 8
    s[8] = makeTableSeries("WC", 1, 3)
    // West Semis: indices 9-10
    s[9]  = makeTableSeries("WS", 1, 4)
    s[10] = makeTableSeries("WL", 2, 2)
    // West R1: indices 11-14
    s[11] = makeTableSeries("W1", 1, 4) // 1v8
    s[12] = makeTableSeries("W4", 4, 3) // 4v5
    s[13] = makeTableSeries("W3", 3, 2) // 3v6
    s[14] = makeTableSeries("W2", 2, 1) // 2v7
    return s
}

func TestBracketTableRenderer_ProducesEightDataLines(t *testing.T) {
    br := newBracketTableRenderer(makeFullBracket(), 0)
    output := br.Render()
    lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
    // header lines (2) + separator (1) + 8 data lines = at least 11
    if len(lines) < 11 {
        t.Errorf("Render() produced %d lines, want at least 11", len(lines))
    }
    // Count non-header lines after separator
    dataStart := 0
    for i, l := range lines {
        if strings.HasPrefix(l, "─") || strings.HasPrefix(strings.TrimSpace(l), "─") {
            dataStart = i + 1
            break
        }
    }
    dataLines := lines[dataStart:]
    if len(dataLines) != 8 {
        t.Errorf("data lines = %d, want 8", len(dataLines))
    }
}

func TestBracketTableRenderer_WestR1AppearsAtLine0(t *testing.T) {
    // West R1 series 0 (index 11) top team should appear on data line 0
    br := newBracketTableRenderer(makeFullBracket(), 0)
    output := br.Render()
    dataLines := extractDataLines(output)
    if !strings.Contains(dataLines[0], "W1T") {
        t.Errorf("data line 0 missing West R1[0] top team W1T: %q", dataLines[0])
    }
}

func TestBracketTableRenderer_WestCFAppearsAtLine2(t *testing.T) {
    // West CF top team (index 8) should appear on data line 2
    br := newBracketTableRenderer(makeFullBracket(), 0)
    output := br.Render()
    dataLines := extractDataLines(output)
    if !strings.Contains(dataLines[2], "WCT") {
        t.Errorf("data line 2 missing West CF top team WCT: %q", dataLines[2])
    }
}

func TestBracketTableRenderer_FinalsAppearsAtLines3And4(t *testing.T) {
    br := newBracketTableRenderer(makeFullBracket(), 0)
    output := br.Render()
    dataLines := extractDataLines(output)
    if !strings.Contains(dataLines[3], "FNT") {
        t.Errorf("data line 3 missing Finals top team FNT: %q", dataLines[3])
    }
    if !strings.Contains(dataLines[4], "FNB") {
        t.Errorf("data line 4 missing Finals bottom team FNB: %q", dataLines[4])
    }
}

func TestBracketTableRenderer_EastR1AppearsOnCorrectLines(t *testing.T) {
    // East R1 series 0 (index 0) top on line 0, bottom on line 1
    br := newBracketTableRenderer(makeFullBracket(), 0)
    output := br.Render()
    dataLines := extractDataLines(output)
    if !strings.Contains(dataLines[0], "E1T") {
        t.Errorf("data line 0 missing East R1[0] top E1T: %q", dataLines[0])
    }
    if !strings.Contains(dataLines[1], "E1B") {
        t.Errorf("data line 1 missing East R1[0] bot E1B: %q", dataLines[1])
    }
}

func TestBracketTableRenderer_CursorNoPanic(t *testing.T) {
    series := makeFullBracket()
    for cursor := 0; cursor < 15; cursor++ {
        br := newBracketTableRenderer(series, cursor)
        _ = br.Render() // must not panic
    }
}

func TestBracketTableRenderer_TBDCellsDimmed(t *testing.T) {
    // TBD cells should contain "TBD"
    s := make([]types.PlayoffSeries, 15)
    for i := range s {
        s[i] = types.PlayoffSeries{
            Status:     "pre",
            TopTeam:    types.PlayoffTeam{Tricode: "TBD"},
            BottomTeam: types.PlayoffTeam{Tricode: "TBD"},
        }
    }
    br := newBracketTableRenderer(s, 0)
    output := br.Render()
    if !strings.Contains(output, "TBD") {
        t.Error("Render() of all-TBD bracket missing TBD text")
    }
}

// extractDataLines skips header and separator lines, returns the 8 data lines.
func extractDataLines(output string) []string {
    lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
    for i, l := range lines {
        if strings.HasPrefix(l, "─") || strings.Contains(l, "────") {
            return lines[i+1:]
        }
    }
    return lines
}
