package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dateChangedMsg struct {
	date string
}

type DateSelector struct {
	date    string
	editing bool
	focused bool
	input   textinput.Model
	err     string
	width   int
}

func NewDateSelector(date string) DateSelector {
	ti := textinput.New()
	ti.Placeholder = "YYYY-MM-DD"
	ti.CharLimit = 10
	ti.Width = 12

	return DateSelector{
		date:    date,
		editing: false,
		focused: true,
		input:   ti,
	}
}

func (ds *DateSelector) SetWidth(w int) {
	ds.width = w
}

func (ds *DateSelector) Focus() {
	ds.focused = true
}

func (ds *DateSelector) Blur() {
	ds.focused = false
	ds.editing = false
}

func (ds *DateSelector) previousDay() {
	t, _ := time.Parse("2006-01-02", ds.date)
	ds.date = t.AddDate(0, 0, -1).Format("2006-01-02")
}

func (ds *DateSelector) nextDay() {
	t, _ := time.Parse("2006-01-02", ds.date)
	next := t.AddDate(0, 0, 1)

	eastern, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(eastern)
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, eastern)

	if !next.After(todayDate) {
		ds.date = next.Format("2006-01-02")
	}
}

func validateDate(input string) error {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	eastern, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(eastern)
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, eastern)

	if t.After(todayDate) {
		return fmt.Errorf("cannot select future dates")
	}
	return nil
}

func (ds DateSelector) Update(msg tea.Msg) (DateSelector, tea.Cmd) {
	if !ds.focused {
		return ds, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		ds.err = ""

		if ds.editing {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				value := ds.input.Value()
				if err := validateDate(value); err != nil {
					ds.err = err.Error()
					return ds, nil
				}
				ds.date = value
				ds.editing = false
				ds.input.Blur()
				return ds, func() tea.Msg { return dateChangedMsg{date: ds.date} }

			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				ds.editing = false
				ds.input.Blur()
				return ds, nil
			}

			var cmd tea.Cmd
			ds.input, cmd = ds.input.Update(msg)
			return ds, cmd
		}

		switch {
		case key.Matches(msg, Keymap.Left):
			ds.previousDay()
			return ds, func() tea.Msg { return dateChangedMsg{date: ds.date} }

		case key.Matches(msg, Keymap.Right):
			old := ds.date
			ds.nextDay()
			if ds.date != old {
				return ds, func() tea.Msg { return dateChangedMsg{date: ds.date} }
			}
			return ds, nil

		case key.Matches(msg, Keymap.Enter):
			ds.editing = true
			ds.input.SetValue(ds.date)
			ds.input.Focus()
			return ds, ds.input.Cursor.BlinkCmd()
		}
	}

	return ds, nil
}

func (ds DateSelector) View() string {
	if ds.editing {
		inputView := ds.input.View()
		if ds.err != "" {
			errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
			inputView += "\n" + errStyle.Render(ds.err)
		}
		return lipgloss.NewStyle().
			Width(ds.width).
			Align(lipgloss.Center).
			Render(inputView)
	}

	dateStyle := lipgloss.NewStyle().Bold(true)
	if ds.focused {
		dateStyle = dateStyle.Foreground(lipgloss.Color("5"))
	}

	arrowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	display := arrowStyle.Render("◀") + "  " + dateStyle.Render(ds.date) + "  " + arrowStyle.Render("▶")

	return lipgloss.NewStyle().
		Width(ds.width).
		Align(lipgloss.Center).
		Render(display)
}
