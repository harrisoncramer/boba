package boba

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// List of all top-level models in the application
type Views map[string]tea.Model

// The Router is responsible for changing the top-level model in the application and triggering any route-based updates with the
type Router struct {
	Model     tea.Model
	Views     Views
	ViewStack []string
	QuitKey   string
}

type NewRouterModelOpts struct {
	View  string
	Views Views
	Quit  string
}

// Creates a new router that is responsible for handling navigation around the application via the changeView function
func NewRouterModel(opts NewRouterModelOpts) tea.Model {
	r := Router{
		Views:   opts.Views,
		QuitKey: opts.Quit,
	}

	r.pushModel(opts.View)
	return r
}

func (m Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := m.handleQuit(msg); cmd != nil { // Our global quit handler shortcuts the event loop
		return m, cmd
	}

	switch msg := msg.(type) {
	// When a component triggers a view change we set the new model
	// and then set router params. This RouterParamsMsg can be detected by components
	// that need query parameters, or other data
	case pushViewMsg:
		m.pushModel(msg.view)
		var cmds []tea.Cmd
		cmds = append(cmds, m.setRouterParams(msg.query), m.Model.Init())
		return m, tea.Sequence(cmds...)
	case replaceViewMsg:
		m.replaceModel(msg.view)
		var cmds []tea.Cmd
		cmds = append(cmds, m.setRouterParams(msg.query), m.Model.Init())
		return m, tea.Sequence(cmds...)
	case popMsg:
		m.pop()
		return m, m.Model.Init()
	}

	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg) // Delegate updates to the model
	return m, cmd
}

func (m Router) View() string {
	base := m.Model.View()       // The model of most of the application
	base += lipgloss.NewStyle(). // Helper text to show the current route
					Foreground(lipgloss.Color("#616161")).
					Render(fmt.Sprintf("\nPath: %s", strings.Join(m.ViewStack, "/")))
	return base
}

func (m Router) Init() tea.Cmd {
	return m.Model.Init()
}

// Handles quit keypresses, shortcuts all other message handling
func (m *Router) handleQuit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.QuitKey:
			return tea.Quit
		}
	}
	return nil
}

// Pops the last view off the stack and navigates to it
func (m *Router) pop() {
	if len(m.ViewStack) < 2 {
		return
	}
	m.ViewStack = m.ViewStack[:len(m.ViewStack)-1]
	m.setModel()
}

// This message is fired when a model calls the Pop method
type popMsg struct{}

// Navigate to the previous top-level model
func Pop() tea.Cmd {
	return func() tea.Msg {
		return popMsg{}
	}
}

// This message is fired when a model calls the Push method
type pushViewMsg struct {
	view  string
	query url.Values
}

// Adds a view to the view stack
func Push(view string) tea.Cmd {
	msg := pushViewMsg{view: view}
	msg.query = parseQuery(view)
	return func() tea.Msg {
		return msg
	}
}

// Takes in a view and gets the appropriate model and pushes it to the router stack
func (m *Router) pushModel(view string) {
	if len(m.ViewStack) == 0 || m.ViewStack[len(m.ViewStack)-1] != view {
		m.ViewStack = append(m.ViewStack, view)
	}
	m.setModel()
}

// This message is fired when a model calls the Replace method
type replaceViewMsg struct {
	view  string
	query url.Values
}

// Replaces the current view in the stack
func Replace(view string) tea.Cmd {
	msg := replaceViewMsg{view: view}
	msg.query = parseQuery(view)
	return func() tea.Msg {
		return msg
	}
}

// Takes in a view and gets the appropriate model and replaces the current one
func (m *Router) replaceModel(view string) {
	m.ViewStack[len(m.ViewStack)-1] = view
	m.setModel()
}

/*
The RouterParamsMsg can be used to pass data to the main route or it's children
by way of a message containing parsed URL values.
E.g some/route?foo=bar or /some/bare/route are both valid.
*/
type RouterParamsMsg struct {
	Params url.Values
}

// Fires when the view is changed. This method fires the
// RouterParamsMsg which can be handled by submodels to get route parameters
func (m *Router) setRouterParams(vals url.Values) tea.Cmd {
	return func() tea.Msg {
		return RouterParamsMsg{vals}
	}
}

// Sets the model in the router based on the last view in the view stack
func (m *Router) setModel() {
	view := m.ViewStack[len(m.ViewStack)-1]
	splitView := strings.Split(view, "?")
	modelName := splitView[0]
	m.Model = m.Views[modelName]
}

// Parses the query values in the view string into url.Values
func parseQuery(view string) url.Values {
	path := strings.Split(view, "?")
	if len(path) == 1 {
		return nil
	}
	query := path[1]
	queryVals, err := url.ParseQuery(query)
	if err != nil {
		log.Fatalf("Error parsing path %v", err)
	}
	return queryVals
}
