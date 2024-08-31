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
		case KeyOpts.Help:
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
	for _, k := range keys {
		switch k {
		case KeyOpts.Back:
			m.Back = key.NewBinding(
				key.WithKeys(KeyOpts.Back),
				key.WithHelp(KeyOpts.Back, "back"),
			)
		case KeyOpts.Quit:
			m.Quit = key.NewBinding(
				key.WithKeys(KeyOpts.Quit),
				key.WithHelp(KeyOpts.Quit, "quit"),
			)
		case KeyOpts.Help:
			m.Help = key.NewBinding(
				key.WithKeys(KeyOpts.Help),
				key.WithHelp(KeyOpts.Help, "help"),
			)
		case KeyOpts.Select:
			m.Select = key.NewBinding(
				key.WithKeys(KeyOpts.Select),
				key.WithHelp(KeyOpts.Select, "select/submit"),
			)
		case KeyOpts.Toggle:
			m.Toggle = key.NewBinding(
				key.WithKeys(KeyOpts.Toggle),
				key.WithHelp(KeyOpts.Toggle, "toggle"),
			)
		case KeyOpts.Up:
			m.Up = key.NewBinding(
				key.WithKeys(KeyOpts.Up),
				key.WithHelp(KeyOpts.Up, "up"),
			)
		case KeyOpts.Down:
			m.Down = key.NewBinding(
				key.WithKeys(KeyOpts.Down),
				key.WithHelp(KeyOpts.Down, "down"),
			)
		case KeyOpts.Filter:
			m.Filter = key.NewBinding(
				key.WithKeys(KeyOpts.Filter),
				key.WithHelp(KeyOpts.Filter, "filter"),
			)
		}
	}
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
