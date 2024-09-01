package boba

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// An individual option in the MultiSelectorModel
type MultiSelectorOption struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Selected bool   `json:"selected"`
	Disabled bool   `json:"disabled"`
}

type MultiSelectorOptions []MultiSelectorOption

// Message that can be used to set all of the options in the model
type MultiSelectorOptionsMsg struct {
	Options MultiSelectorOptions
}

type MultiSelectorModel struct {
	cursor         int
	cursorIcon     string
	options        MultiSelectorOptions
	visibleOptions MultiSelectorOptions
	filter         textinput.Model
	keys           KeyOpts
	theme          Theme
	name           string
	maxHeight      func() int
	truncated      bool
	LoadingModel
}

type NewMultiSelectorModelOpts struct {
	Filter    FilterOpts
	Options   []MultiSelectorOption
	Theme     Theme
	Name      string
	MaxHeight func() int
	Keys      KeyOpts
}

// Allows for the toggling of multiple values in a list via a toggle mechanism.
// The selection triggers the
func NewMultiSelectorModel(opts NewMultiSelectorModelOpts) MultiSelectorModel {
	m := MultiSelectorModel{
		options:        opts.Options,
		visibleOptions: opts.Options,
		theme:          opts.Theme,
		maxHeight:      opts.MaxHeight,
		keys:           opts.Keys,
		LoadingModel:   NewLoadingModel(),
	}

	if !opts.Filter.Hidden {
		ti := textinput.New()
		ti.Placeholder = opts.Filter.Placeholder
		m.filter = ti
	}

	return m
}

// The unselectAllMsg is used to reset the state of the component
type unselectAllMsg struct{}

func (m MultiSelectorModel) Init() tea.Cmd {
	return func() tea.Msg {
		return unselectAllMsg{}
	}
}

// Message used to set options in a MultiSelectorModel
func (m MultiSelectorModel) Update(msg tea.Msg) (MultiSelectorModel, tea.Cmd) {
	var cmds []tea.Cmd

	m.filter = UpdateSubmodel(m.filter, msg, &cmds)
	m.Spinner = m.UpdateLoading(msg, &cmds)

	switch msg := msg.(type) {
	case unselectAllMsg:
		m.unselectAll()
	case MultiSelectorOptionsMsg:
		m.setOptions(msg.Options)
	case tea.KeyMsg:
		switch msg.String() {
		case m.keys.Down:
			m.move(Down)
		case m.keys.Up:
			m.move(Up)
		case m.keys.Toggle:
			cmds = append(cmds, m.toggleVal)
		case m.keys.Filter:
			cmds = append(cmds, textinput.Blink)
			m.filter.Focus()
		case m.keys.Back:
			if m.filter.Focused() {
				m.filter.Blur()
				return m, nil
			} else {
				return m, back(m.name)
			}
		}
	}

	if m.filter.Focused() {
		m.cursor = 0
	}

	m.filterOptions()

	return m, tea.Batch(cmds...)
}

func (m MultiSelectorModel) View() string {
	if m.Loading {
		return fmt.Sprintf("\n%s\n", m.Spinner.View())
	}
	base := strings.Builder{}
	base.WriteString(rebuildCursor(m.filter.View(), m.filter.Focused(), m.theme))
	if len(m.visibleOptions) == 0 {
		base.WriteString("No options found \n")
	} else {
		for i, option := range m.visibleOptions {
			icon := "[x]"
			if !option.Selected {
				icon = "[ ]"
			}

			var color ColorType
			if option.Disabled {
				color = Neutral
			}

			if i == m.cursor {
				icon = m.theme.Color(icon, color)
				base.WriteString(fmt.Sprintf("%s %s ", m.theme.ColorCond(">", Primary, !m.FilterFocused()), icon))
			} else {
				icon = m.theme.Color(icon, color)
				base.WriteString(fmt.Sprintf("%s  %s ", strings.Repeat(" ", len(m.cursorIcon)), icon))
			}

			base.WriteString(fmt.Sprintf("%s\n", m.theme.Color(option.Label, color)))
		}

		if m.truncated {
			base.WriteString(m.theme.Color(fmt.Sprintf("  Results limited, use %s to search...\n", m.keys.Filter), Neutral))
		}
	}
	return base.String()
}

// Indicates whether the multi-selector is in a focused state
func (m MultiSelectorModel) FilterFocused() bool {
	return m.filter.Focused()
}

// Moves the cursor up or down among the options
func (m *MultiSelectorModel) move(direction Direction) {
	if m.filter.Focused() {
		return
	}
	if direction == Up {
		if m.cursor > 0 {
			m.cursor--
		}
	} else {
		if m.cursor < len(m.visibleOptions)-1 {
			m.cursor++
		}
	}
}

// Removes selection status from all possible options prior to quit
func (m *MultiSelectorModel) unselectAll() {
	var results []MultiSelectorOption
	for _, opt := range m.options {
		opt.Selected = false
		results = append(results, opt)
	}
	m.options = results
}

// Sets options on the selector
func (m *MultiSelectorModel) setOptions(options []MultiSelectorOption) {
	m.options = options
}

// Filters the possible options by the text contained in the textinput model
func (m *MultiSelectorModel) filterOptions() {
	var visibleOptions MultiSelectorOptions
	for _, opt := range m.options {
		filter := m.filter.Value()
		if filter == "" || strings.Contains(strings.ToLower(opt.Label), strings.ToLower(filter)) {
			visibleOptions = append(visibleOptions, opt)
		}
	}

	// If we have exceeded the max height, trim our results
	if m.maxHeight != nil {
		h := m.maxHeight() - 2 // Include the height of the filter input
		if len(visibleOptions) > h {
			visibleOptions = visibleOptions[:h]
			m.truncated = true
		} else {
			m.truncated = false
		}
	}

	m.visibleOptions = visibleOptions
}

// Message for a single toggle event
type MultiSelectorOptionMsg struct {
	Option MultiSelectorOption
}

// Toggles the current value. Triggers the MultiSelectorOptionMsg which can be handled
// by other components
func (m *MultiSelectorModel) toggleVal() tea.Msg {
	if !m.FilterFocused() {
		val := m.visibleOptions[m.cursor].Value
		i := findIndex(m.options, func(opt MultiSelectorOption) bool {
			return opt.Value == val
		})
		if m.options[i].Disabled {
			return nil
		}
		m.options[i].Selected = !m.options[i].Selected
		return MultiSelectorOptionMsg{m.options[i]}
	}
	return nil
}
