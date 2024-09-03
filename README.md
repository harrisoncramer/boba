# Boba 🧋

This is my component library for BubbleTea.

It is a work in progress. Please do not use this library in your own projects (yet)!

## Component Usage

Install the package:

```bash
go get github.com/harrisoncramer/boba
```

Then use the components. For instance, here's how you can use the selector component, which allows you to choose from one of several options.

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/harrisoncramer/boba"
)

type MyModel struct {
	selector boba.SelectorModel
}

func NewModel() tea.Model {
	return MyModel{
		selector: boba.NewSelectorModel(boba.NewSelectorModelOpts{
			Keys: boba.KeyOpts{
				Back:   "esc",
				Up:     "up",
				Down:   "down",
				Select: "enter",
				Filter: "/",
			},
			Filter: boba.FilterOpts{
				Placeholder: "Search...",
			},
			Options: []boba.SelectorOption{
				{Label: "Value One", Value: "value_1"},
				{Label: "Value Two", Value: "value_2"},
				{Label: "Value Three", Value: "value_3"},
				{Label: "Value Four", Value: "value_4"},
				{Label: "Value Five", Value: "value_5"},
			},
		}),
	}
}

func (m MyModel) Init() tea.Cmd {
	return nil
}

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case boba.SelectMsg:
		fmt.Printf("You chose %s", msg.Option.Value)
		os.Exit(0)
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	m.selector, cmd = m.selector.Update(msg)
	return m, cmd
}

func (m MyModel) View() string {
	base := "My Program\n\n"
	base += m.selector.View()
	return base
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
```

## Router Usage

You can set up a router with various views and child routes like this:

```go
	router := router.NewRouterModel(router.NewRouterModelOpts{
		View:        "root",
		DefaultView: "root",
		Quit:        shared.PluginOptions.Keys.Quit,
		Views: router.Views{
			{
				Path: "root",
				Model: NewMainModel(),
				Children: router.Views{
					{
						Path: "subview-1",
						Model: NewSubview1(),
					},
					{
						Path: "subview-2",
						Model: NewSubview2(),
						Children: router.Views{
							{
								Path: "subview-3?foo=bar",
								Model: NewSubview3()
							},
							{
								Path: "subview-4",
								Model: NewSubview4()
							},
						},
                    },
				},
			},
		},
	})
```

And then to navigate around, you can use the methods exposed from the router package such as:

```go
router.Push("subview-1")
router.Pop() // Navigates back to root
router.Push("subview-3") // Directly navigate to subview 3, bypass parent
router.GetParam("foo") // "bar"
```

This router does not support URL parameters (like :id), instead just encode that data into query parameters. Route matching is done on simple strings, not regular expressions.
