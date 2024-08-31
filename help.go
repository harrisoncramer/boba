package boba

import (
	"fmt"
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
	kM := NewKeys(keys...)
	return HelpModel{
		help: help.New(),
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
	base := strings.Builder{}
	base.WriteString(fmt.Sprintf("\n%s", m.help.View(m.keys)))

	return base.String()
}

type keyMap struct {
	Quit   key.Binding
	Back   key.Binding
	Select key.Binding
	Toggle key.Binding
	Up     key.Binding
	Down   key.Binding
	Filter key.Binding
	Help   key.Binding
}

func NewKeys(keys ...string) keyMap {
	var m keyMap
	// for _, k := range keys {
	// 	switch k {
	// 	case m.Back:
	// 		m.Back = key.NewBinding(
	// 			key.WithKeys(m.keys.Back),
	// 			key.WithHelp(m.keys.Back, "back"),
	// 		)
	// 	case m.keys.Quit:
	// 		m.Quit = key.NewBinding(
	// 			key.WithKeys(m.keys.Quit),
	// 			key.WithHelp(m.keys.Quit, "quit"),
	// 		)
	// 	case m.keys.Help:
	// 		m.Help = key.NewBinding(
	// 			key.WithKeys(m.keys.Help),
	// 			key.WithHelp(m.keys.Help, "help"),
	// 		)
	// 	case m.keys.Select:
	// 		m.Select = key.NewBinding(
	// 			key.WithKeys(m.keys.Select),
	// 			key.WithHelp(m.keys.Select, "select/submit"),
	// 		)
	// 	case m.keys.Toggle:
	// 		m.Toggle = key.NewBinding(
	// 			key.WithKeys(m.keys.Toggle),
	// 			key.WithHelp(m.keys.Toggle, "toggle"),
	// 		)
	// 	case m.keys.Up:
	// 		m.Up = key.NewBinding(
	// 			key.WithKeys(m.keys.Up),
	// 			key.WithHelp(m.keys.Up, "up"),
	// 		)
	// 	case m.keys.Down:
	// 		m.Down = key.NewBinding(
	// 			key.WithKeys(m.keys.Down),
	// 			key.WithHelp(m.keys.Down, "down"),
	// 		)
	// 	case m.keys.Filter:
	// 		m.Filter = key.NewBinding(
	// 			key.WithKeys(m.keys.Filter),
	// 			key.WithHelp(m.keys.Filter, "filter"),
	// 		)
	// 	}
	// }
	return m
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Back,
		},
		{
			k.Quit,
		},
		{
			k.Help,
		},
		{
			k.Select,
		},
		{
			k.Toggle,
		},
		{
			k.Up,
		},
		{
			k.Down,
		},
		{
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
