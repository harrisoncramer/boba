package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/harrisoncramer/cbcli/shared"
)

type TextInputModel struct {
	theme  Theme
	id     string
	input  textinput.Model
	noUp   bool
	noDown bool
}

type NewTextInputOptions struct {
	NoDown      bool
	NoUp        bool
	Id          string
	Placeholder string
	Theme       Theme
}

// Wrapper around the textinput model from BubbleTea. Extended to handle focusing and
// moving up/down to the nearest ComponentModel
func NewTextInputModel(opts NewTextInputOptions, models ...textinput.Model) ComponentModel {
	ti := TextInputModel{
		input:  textinput.New(),
		id:     opts.Id,
		noUp:   opts.NoUp,
		noDown: opts.NoDown,
		theme:  opts.Theme,
	}
	ti.input.Placeholder = opts.Placeholder
	return &ti
}

func (m TextInputModel) Init() tea.Cmd {
	return m.input.Focus()
}

func (m TextInputModel) Update(msg tea.Msg) (ComponentModel, tea.Cmd) {
	var cmds = []tea.Cmd{}

	m.input = shared.UpdateSubmodel(m.input, msg, &cmds)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case shared.PluginOptions.Keys.Up:
			if m.Focused() && !m.noUp {
				m.Blur()
				return &m, back(m.id)
			}
		case shared.PluginOptions.Keys.Down:
			if m.Focused() && !m.noDown {
				m.Blur()
				return &m, next(m.id)
			}
		case shared.PluginOptions.Keys.Back:
			m.Blur()
			return &m, nil
		}
	}

	return &m, tea.Batch(cmds...)
}

func (m TextInputModel) View() string {
	return rebuildCursor(m.input.View(), m.input.Focused(), m.theme)
}

func (m *TextInputModel) Blur() {
	m.input.Blur()
}

func (m *TextInputModel) Clear() {
	m.input.SetValue("")
}

func (m *TextInputModel) Focus() tea.Cmd {
	return m.input.Focus()
}

func (m TextInputModel) Focused() bool {
	return m.input.Focused()
}

func (m TextInputModel) Id() string {
	return m.id
}

func (m TextInputModel) Value() any {
	return m.input.Value()
}
