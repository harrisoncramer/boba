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

// The stack of views in the router
var viewStack []string

type Router struct {
	Model   tea.Model
	Views   Views
	QuitKey string
}

type NewRouterModelOpts struct {
	View  string
	Views Views
	Quit  string
}

// The Router is responsible for changing the top-level model in the application and triggering any route-based updates
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

	// When a component triggers a view change we set the new model
	// and then set router params. This RouterParamsMsg can be detected by components
	// that need query parameters, or other data
	switch msg := msg.(type) {
	case pushViewMsg:
		m.pushModel(msg.view)
		return m, m.Model.Init()
	case replaceViewMsg:
		m.replaceModel(msg.view)
		return m, m.Model.Init()
	case popMsg:
		m.popModel()
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
					Render(fmt.Sprintf("\nPath: %s", strings.Join(viewStack, "/")))
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
func (m *Router) popModel() {
	if len(viewStack) < 2 {
		return
	}
	viewStack = viewStack[:len(viewStack)-1]
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
	if len(viewStack) == 0 || viewStack[len(viewStack)-1] != view {
		viewStack = append(viewStack, view)
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
	viewStack[len(viewStack)-1] = view
	m.setModel()
}

// Provides the current url values parsed from the top route. Can be called
// by components to get the params
func GetParams() url.Values {
	return parseQuery(viewStack[len(viewStack)-1])
}

// Provides a specific url parameter from the top route
func GetParam(val string) string {
	queries := parseQuery(viewStack[len(viewStack)-1])
	return queries.Get(val)
}

// The message triggered when a router parameter has changed
type RouterParamChangedMsg struct {
	Key   string
	Value string
}

// Provides the ability to change url parameters on the fly within components
func SetParam(key string, val string) tea.Cmd {
	return func() tea.Msg {
		view := viewStack[len(viewStack)-1]
		query := parseQuery(view)
		if query == nil {
			viewStack[len(viewStack)-1] = fmt.Sprintf("%s?%s=%s", view, key, val)
			return nil
		}
		query.Set(key, val)
		log.Printf("Query is %s", query.Encode())
		newQuery := query.Encode()
		baseView := strings.Split(view, "?")[0] // TODO: Nested views with query params
		viewStack[len(viewStack)-1] = fmt.Sprintf("%s?%s", baseView, newQuery)
		return RouterParamChangedMsg{
			Key:   key,
			Value: val,
		}
	}
}

// Sets the model in the router based on the last view in the view stack
func (m *Router) setModel() {
	view := viewStack[len(viewStack)-1]
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
