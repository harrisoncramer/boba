package boba

import "github.com/charmbracelet/bubbles/key"

// The KeyOpts struct contains all possible keys and can be used
// to load in a set of keys for boba
type KeyOpts struct {
	Up     string `mapstructure:"up"`
	Down   string `mapstructure:"down"`
	Select string `mapstructure:"select"`
	Toggle string `mapstructure:"toggle"`
	Back   string `mapstructure:"back"`
	Quit   string `mapstructure:"quit"`
	Filter string `mapstructure:"filter"`
	Help   string `mapstructure:"help"`
}

// Contains a mapping of all keys to their key bindings (bubbletea type)
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

var bobaKeys KeyOpts

// Sets the keys in the boba components
func SetKeys(k KeyOpts) {
	bobaKeys = k
}

// Creates a key map that's a subset of available keys, for use in different
// components that only use some of them
func newKeys(keys ...string) keyMap {
	var m keyMap
	for _, k := range keys {
		switch k {
		case bobaKeys.Quit:
			m.Quit = key.NewBinding(
				key.WithKeys(bobaKeys.Quit),
				key.WithHelp(bobaKeys.Quit, "quit"),
			)
		case bobaKeys.Select:
			m.Select = key.NewBinding(
				key.WithKeys(bobaKeys.Select),
				key.WithHelp(bobaKeys.Select, "select/submit"),
			)
		case bobaKeys.Toggle:
			m.Toggle = key.NewBinding(
				key.WithKeys(bobaKeys.Toggle),
				key.WithHelp(bobaKeys.Toggle, "toggle"),
			)
		case bobaKeys.Up:
			m.Up = key.NewBinding(
				key.WithKeys(bobaKeys.Up),
				key.WithHelp(bobaKeys.Up, "up"),
			)
		case bobaKeys.Down:
			m.Down = key.NewBinding(
				key.WithKeys(bobaKeys.Down),
				key.WithHelp(bobaKeys.Down, "down"),
			)
		case bobaKeys.Back:
			m.Back = key.NewBinding(
				key.WithKeys(bobaKeys.Back),
				key.WithHelp(bobaKeys.Back, "back"),
			)
		case bobaKeys.Filter:
			m.Filter = key.NewBinding(
				key.WithKeys(bobaKeys.Filter),
				key.WithHelp(bobaKeys.Filter, "filter"),
			)
		case bobaKeys.Help:
			m.Help = key.NewBinding(
				key.WithKeys(bobaKeys.Help),
				key.WithHelp(bobaKeys.Help, "help"),
			)
		}
	}
	return m
}
