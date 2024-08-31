package boba

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type HelpModel struct {
	keys keyMap
	help help.Model
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func NewHelpModel(keys ...string) HelpModel {
	kM := newKeys(keys...)
	h := help.New()
	h.ShortSeparator = " • "
	h.FullSeparator = " "
	return HelpModel{
		help: h,
		keys: kM,
	}
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.keys.Help.Help().Key:
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

func (m HelpModel) View() string {
	base := m.help.View(m.keys)
	base = strings.ReplaceAll(base, "\n", "")
	return "\n" + strings.Join(strings.Fields(base), " ")
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Back,
			k.Quit,
			k.Help,
			k.Select,
			k.Toggle,
			k.Up,
			k.Down,
			k.Filter,
		},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Back,
		k.Quit,
		k.Help,
	}
}
