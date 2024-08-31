package boba

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorType string

type Colors struct {
	Success   string `mapstructure:"success"`
	Neutral   string `mapstructure:"neutral"`
	Primary   string `mapstructure:"primary"`
	Secondary string `mapstructure:"secondary"`
}

// Possible types of colors
const (
	Success   ColorType = "Success"
	Primary   ColorType = "Primary"
	Neutral   ColorType = "Neutral"
	Secondary ColorType = "Secondary"
)

// Used to set colors and styles, e.g. t.color("some-text", Success)
type Theme map[ColorType]lipgloss.Style

// Creates a new theme for the components with overrideable defaults
func NewTheme(overrides Colors) Theme {
	defaultColors := map[ColorType]string{
		Primary:   "#78A7D8",
		Secondary: "#FFA066",
		Neutral:   "#979797",
		Success:   "#98BB6C",
	}
	if overrides.Primary != "" {
		defaultColors[Primary] = overrides.Primary
	}
	if overrides.Secondary != "" {
		defaultColors[Secondary] = overrides.Secondary
	}
	if overrides.Neutral != "" {
		defaultColors[Neutral] = overrides.Neutral
	}
	if overrides.Success != "" {
		defaultColors[Success] = overrides.Success
	}
	t := make(Theme)
	for key, color := range defaultColors {
		t[key] = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	}
	return t
}

// Applies the relevant color from the theme to the text
func (t Theme) Color(text string, ct ColorType) string {
	return t[ct].Render(text)
}

// Applies the color to the text, if the provided condition is true
func (t Theme) ColorCond(text string, ct ColorType, condition bool) string {
	if !condition {
		return text
	}
	return t[ct].Render(text)
}
