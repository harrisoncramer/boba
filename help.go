package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/harrisoncramer/cbcli/shared"
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
		case shared.PluginOptions.Keys.Help:
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
		case shared.PluginOptions.Keys.Back:
			m.Back = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Back),
				key.WithHelp(shared.PluginOptions.Keys.Back, "back"),
			)
		case shared.PluginOptions.Keys.Quit:
			m.Quit = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Quit),
				key.WithHelp(shared.PluginOptions.Keys.Quit, "quit"),
			)
		case shared.PluginOptions.Keys.Help:
			m.Help = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Help),
				key.WithHelp(shared.PluginOptions.Keys.Help, "help"),
			)
		case shared.PluginOptions.Keys.Select:
			m.Select = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Select),
				key.WithHelp(shared.PluginOptions.Keys.Select, "select/submit"),
			)
		case shared.PluginOptions.Keys.Toggle:
			m.Toggle = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Toggle),
				key.WithHelp(shared.PluginOptions.Keys.Toggle, "toggle"),
			)
		case shared.PluginOptions.Keys.Up:
			m.Up = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Up),
				key.WithHelp(shared.PluginOptions.Keys.Up, "up"),
			)
		case shared.PluginOptions.Keys.Down:
			m.Down = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Down),
				key.WithHelp(shared.PluginOptions.Keys.Down, "down"),
			)
		case shared.PluginOptions.Keys.Filter:
			m.Filter = key.NewBinding(
				key.WithKeys(shared.PluginOptions.Keys.Filter),
				key.WithHelp(shared.PluginOptions.Keys.Filter, "filter"),
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
