package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ComponentBackMsg struct {
	ComponentName string
}

type ComponentNextMsg struct {
	ComponentName string
}

type ComponentFocusMsg string

func (f ComponentFocusMsg) String() string {
	return string(f)
}

func back(componentName string) tea.Cmd {
	return func() tea.Msg {
		return ComponentBackMsg{componentName}
	}
}

func next(componentName string) tea.Cmd {
	return func() tea.Msg {
		return ComponentNextMsg{componentName}
	}
}

type ComponentModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (ComponentModel, tea.Cmd)
	View() string
	Clear()
	Focus() tea.Cmd
	Focused() bool
	Value() any
	Blur()
	Id() string
}

type Direction string

const (
	Up   Direction = "up"
	Down Direction = "down"
)

type FilterOpts struct {
	Placeholder string
	Hidden      bool
}

// Used to re-color the cursor that bubbletea provides, ugh
func rebuildCursor(rawString string, focused bool, theme Theme) string {
	base := strings.Builder{}
	if focused {
		_, after, _ := strings.Cut(rawString, "> ")
		base.WriteString(theme.Color("> ", Primary))
		base.WriteString(after)
	} else {
		base.WriteString(rawString)
	}

	base.WriteString("\n")

	return base.String()
}
