package boba

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ToggleModel struct {
	name    string
	on      bool
	label   string
	focused bool
	noUp    bool
	noDown  bool
	theme   Theme
	keys    KeyOpts
}

type NewToggleOptions struct {
	Label  string
	Name   string
	NoDown bool
	NoUp   bool
	Theme  Theme
	Keys   KeyOpts
	On     bool
}

type SetToggleMsg struct{ On bool }

// Allows for turning a boolean value true/false, and adheres to the ComponentModel
// to be used in forms
func NewToggleModel(opts NewToggleOptions) ComponentModel {
	return &ToggleModel{
		on:     opts.On,
		label:  opts.Label,
		name:   opts.Name,
		noDown: opts.NoDown,
		noUp:   opts.NoUp,
		theme:  opts.Theme,
		keys:   opts.Keys,
	}
}

func (m ToggleModel) Init() tea.Cmd {
	return nil
}

func (m ToggleModel) Update(msg tea.Msg) (ComponentModel, tea.Cmd) {
	if !m.Focused() {
		return &m, nil
	}
	switch msg := msg.(type) {
	case SetToggleMsg:
		m.on = msg.On
	case tea.KeyMsg:
		switch msg.String() {
		case m.keys.Toggle:
			m.on = !m.on
			return &m, m.changeToggle
		case m.keys.Up:
			if m.Focused() && !m.noUp {
				m.Blur()
				return &m, back(m.name)
			}
		case m.keys.Down:
			if m.Focused() && !m.noDown {
				m.Blur()
				return &m, next(m.name)
			}
		}
	}

	return &m, nil
}

type ChangeToggleMsg struct{ On bool }

func (m ToggleModel) changeToggle() tea.Msg {
	return ChangeToggleMsg{On: m.on}
}

func (m ToggleModel) View() string {
	base := fmt.Sprintf("%s %s: ", m.theme.ColorCond(">", Primary, m.Focused()), m.label)
	if m.on {
		base += m.theme.Color("Yes", Success)
	} else {
		base += m.theme.Color("No", Neutral)
	}

	base += "\n"
	return base
}

func (m ToggleModel) Focused() bool {
	return m.focused
}

func (m *ToggleModel) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *ToggleModel) Blur() {
	m.focused = false
}

func (m ToggleModel) Id() string {
	return m.name
}

func (m *ToggleModel) Clear() {
	m.focused = false
	m.on = false
}

func (m ToggleModel) Value() any {
	return m.on
}
