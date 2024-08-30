package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/harrisoncramer/cbcli/shared"
)

type ToggleModel struct {
	name    string
	on      bool
	label   string
	focused bool
	noUp    bool
	noDown  bool
	theme   Theme
}

type NewToggleOptions struct {
	Label  string
	Name   string
	NoDown bool
	NoUp   bool
	Theme  Theme
}

// Allows for turning a boolean value true/false, and adheres to the ComponentModel
// to be used in forms
func NewToggleModel(opts NewToggleOptions) ComponentModel {
	return &ToggleModel{
		on:     false,
		label:  opts.Label,
		name:   opts.Name,
		noDown: opts.NoDown,
		noUp:   opts.NoUp,
		theme:  opts.Theme,
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
	case tea.KeyMsg:
		switch msg.String() {
		case shared.PluginOptions.Keys.Toggle:
			m.on = !m.on
		case shared.PluginOptions.Keys.Up:
			if m.Focused() && !m.noUp {
				m.Blur()
				return &m, back(m.name)
			}
		case shared.PluginOptions.Keys.Down:
			if m.Focused() && !m.noDown {
				m.Blur()
				return &m, next(m.name)
			}
		}
	}

	return &m, nil
}

func (m ToggleModel) View() string {
	base := fmt.Sprintf("%s %s: ", m.theme.ColorCond(">", Primary, m.Focused()), m.label)
	if m.on {
		base += m.theme.Color("Yes", Success)
	} else {
		base += m.theme.Color("No", Danger)
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
