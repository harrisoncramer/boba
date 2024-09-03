package boba

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	Path     string
	Model    tea.Model
	Children []View
}

func (v View) isMatch(path string) bool {
	view := strings.Split(path, "?")[0] // TODO: Nested views with query params
	return v.Path == view
}

// List of all top-level models in the application
type Views []View

// The stack of views in the router
var viewStack []string

type Router struct {
	Model       tea.Model
	Views       Views
	DefaultView string // View that is navigated to when "back" is called w/out a previous route
	QuitKey     string
}

type NewRouterModelOpts struct {
	View        string
	Views       Views
	Quit        string
	DefaultView string
}

// The Router is responsible for changing the top-level model in the application and triggering any route-based updates
// Creates a new router that is responsible for handling navigation around the application via the changeView function
func NewRouterModel(opts NewRouterModelOpts) tea.Model {
	r := Router{
		Views:       opts.Views,
		DefaultView: opts.DefaultView,
		QuitKey:     opts.Quit,
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
	base := m.Model.View() // The model of most of the application

	currentView := viewStack[len(viewStack)-1]
	base += lipgloss.NewStyle(). // Helper text to show the current route
					Foreground(lipgloss.Color("#616161")).
					Render(fmt.Sprintf("\nPath: %s", currentView))
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
	if len(viewStack) < 2 && m.DefaultView == "" {
		return
	}

	// If a default view is provided and there is no previous view, then navigate to it
	if len(viewStack) < 2 {
		viewStack = []string{m.DefaultView}
		m.setModel(m.DefaultView, m.Views)
		return
	}

	viewStack = viewStack[:len(viewStack)-1]
	m.setModel(viewStack[len(viewStack)-1], m.Views)
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
	viewStack = append(viewStack, view)
	m.setModel(view, m.Views)
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
	m.setModel(view, m.Views)
}

// Provides the current url values parsed from the top route. Can be called
// by components to get the params
func GetQueryParams() url.Values {
	return parseQuery(viewStack[len(viewStack)-1])
}

// Provides a specific url parameter from the top route
func GetQueryParam(val string) string {
	queries := parseQuery(viewStack[len(viewStack)-1])
	return queries.Get(val)
}

// Returns the current route
func GetRoute() string {
	return viewStack[len(viewStack)-1]
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
		queryString := parseQuery(view)
		if queryString == nil {
			viewStack[len(viewStack)-1] = fmt.Sprintf("%s?%s=%s", view, key, val)
			return nil
		}
		queryString.Set(key, val)
		newQuery := queryString.Encode()
		baseView := strings.Split(view, "?")[0] // TODO: Nested views with query params
		viewStack[len(viewStack)-1] = fmt.Sprintf("%s?%s", baseView, newQuery)
		return RouterParamChangedMsg{
			Key:   key,
			Value: val,
		}
	}
}

// Depth first search of the router until it finds the matching route
func (m *Router) setModel(view string, views Views) bool {
	for _, v := range views {
		if m.setModel(view, v.Children) {
			return true
		}
		if v.isMatch(view) {
			m.Model = v.Model
			return true
		}
	}
	return false
}

// Parses the query values in the view string into url.Values
func parseQuery(view string) url.Values {
	queryString := getLastQueryString(view)
	if queryString == "" {
		return nil
	}
	queryVals, err := url.ParseQuery(queryString)
	if err != nil {
		log.Fatalf("Error parsing path %v", err)
	}

	return queryVals
}

func getLastQueryString(view string) string {
	idx := strings.LastIndex(view, "?")
	if idx == -1 {
		return ""
	}
	return view[idx+1:]
}
