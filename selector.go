package boba

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectorOption struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

type SelectorOptions []SelectorOption

// Displays a list of fitlerable and selectable options
type SelectorModel struct {
	cursor         int
	cursorIcon     string
	options        SelectorOptions
	visibleOptions SelectorOptions
	filter         textinput.Model
	theme          Theme
	name           string
	maxHeight      func() int
	truncated      bool
	keys           KeyOpts
	LoadingModel
}

type NewSelectorModelOpts struct {
	Filter    FilterOpts
	Options   []SelectorOption
	Theme     Theme
	Name      string
	MaxHeight func() int
	Keys      KeyOpts
}

// Allows for the selection of a single value among a list of options.
// The selection triggers the SelectMsg which can be handled by other models.
func NewSelectorModel(opts NewSelectorModelOpts) SelectorModel {
	m := SelectorModel{
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

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

// Used to set the options in the model
type SelectorOptionsMsg struct {
	Options SelectorOptions
}

func (m SelectorModel) Update(msg tea.Msg) (SelectorModel, tea.Cmd) {
	var cmds []tea.Cmd

	m.filter = UpdateSubmodel(m.filter, msg, &cmds)
	m.spinner = m.updateLoading(msg, &cmds)

	switch msg := msg.(type) {
	case SelectorOptionsMsg:
		m.setOptions(msg.Options)
	case tea.KeyMsg:
		switch msg.String() {
		case m.keys.Down:
			m.move(Down)
		case m.keys.Up:
			m.move(Up)
		case m.keys.Select:
			if !m.filter.Focused() {
				return m, m.selectVal
			} else {
				m.filter.Blur()
			}
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

	// Reset the cursor when someone filters the list
	if m.filter.Focused() {
		m.cursor = 0
	}

	m.filterOptions() // Use the filter to update the list of options

	return m, tea.Batch(cmds...)
}

func (m SelectorModel) View() string {
	if m.loading {
		return fmt.Sprintf("\n%s\n", m.spinner.View())
	}
	base := strings.Builder{}
	base.WriteString(rebuildCursor(m.filter.View(), m.filter.Focused(), m.theme))
	if len(m.visibleOptions) == 0 {
		base.WriteString("No options found \n")
	} else {
		for i, option := range m.visibleOptions {
			if i == m.cursor {
				base.WriteString(fmt.Sprintf("%s ", m.theme.ColorCond(">", Primary, !m.filter.Focused())))
			} else {
				base.WriteString(fmt.Sprintf("%s  ", strings.Repeat(" ", len(m.cursorIcon))))
			}

			var color ColorType
			if option.Disabled {
				color = Neutral
			}

			base.WriteString(fmt.Sprintf("%s\n", m.theme.Color(option.Label, color)))
		}
		if m.truncated {
			base.WriteString(m.theme.Color(fmt.Sprintf("  Results limited, use %s to search...\n", m.keys.Filter), Neutral))
		}
	}

	return base.String()
}

// Filters the possible options by the text contained in the textinput model
func (m *SelectorModel) filterOptions() {
	var visibleOptions SelectorOptions
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

// Moves the cursor up or down among the options
func (m *SelectorModel) move(direction Direction) {
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

// Sets options on the selector, triggered via the SelectorOptionsMsg
func (m *SelectorModel) setOptions(options []SelectorOption) {
	m.options = options
}

type SelectMsg struct {
	Option SelectorOption
}

// Chooses the value at the given index, triggers the SelectMsg for use in other models
func (m *SelectorModel) selectVal() tea.Msg {
	if !m.filter.Focused() {
		val := m.visibleOptions[m.cursor].Value
		i := findIndex(m.options, func(opt SelectorOption) bool {
			return opt.Value == val
		})
		if m.options[i].Disabled {
			return nil
		}
		return SelectMsg{m.options[i]}
	}
	return nil
}
