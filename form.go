package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Form is a utility for wrapping multiple components which allows for cycling
// the different inputs more easily. Components that send the ComponentBackMsg and
// ComponentNextMsg will trigger focus to shift to the next/last component in the
// list.
type Form []ComponentModel

type StartMsg struct{}

// Navigates to the previous component in the form, if possible
func (f Form) getPrevComponent(msg ComponentBackMsg) ComponentModel {
	i := findIndex(f, func(c ComponentModel) bool {
		return c.Id() == msg.ComponentName
	})
	return f[i-1]
}

// Navigates to the next component in the form, if possible
func (f Form) getNextComponent(msg ComponentNextMsg) ComponentModel {
	i := findIndex(f, func(c ComponentModel) bool {
		return c.Id() == msg.ComponentName
	})
	return f[i+1]
}

// Focuses on the first component in the form
func (f Form) focus() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range f[0:] {
		c.Blur()
		c.Clear()
	}
	cmds = append(cmds, f[0].Focus())
	return tea.Batch(cmds...)
}

func (f Form) Update(msg tea.Msg) (Form, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case ComponentBackMsg:
		c := f.getPrevComponent(msg)
		cmds = append(cmds, c.Focus())
	case ComponentNextMsg:
		c := f.getNextComponent(msg)
		cmds = append(cmds, c.Focus())
	case StartMsg:
		cmds = append(cmds, f.focus()) // The StartMsg triggers focus in the form
	default:
		for i, c := range f {
			var cmd tea.Cmd
			f[i], cmd = c.Update(msg) // Delegate all other messages to the child models
			cmds = append(cmds, cmd)
		}
	}
	return f, tea.Batch(cmds...)
}

func findIndex[T any](slice []T, predicate func(T) bool) int {
	for i, v := range slice {
		if predicate(v) {
			return i
		}
	}
	return -1 // Return -1 if no element satisfies the condition
}
