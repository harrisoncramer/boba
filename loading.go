package boba

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func NewLoadingModel() LoadingModel {
	return LoadingModel{
		Loading: false,
		Loader:  spinner.New(),
	}
}

type LoadingModel struct {
	Loading bool
	Loader  spinner.Model
}

func (loader *LoadingModel) Load() tea.Msg {
	if !loader.Loading {
		return loadingMsg{}
	}
	return nil
}

type SuccessMsg struct{ Msg string }
type ErrMsg struct{ Err error }

func (loader *LoadingModel) UpdateLoading(msg tea.Msg, cmds *[]tea.Cmd) spinner.Model {
	switch msg := msg.(type) {
	case loadingMsg:
		loader.Loading = true
		*cmds = append(*cmds, loader.Loader.Tick)
	case spinner.TickMsg:
		if loader.Loading {
			loader.Loader = UpdateSubmodel(loader.Loader, msg, cmds)
		}
	case MultiSelectorOptionsMsg, SelectorOptionsMsg, SuccessMsg, ErrMsg:
		loader.Loading = false
	}

	return loader.Loader
}

func (m LoadingModel) View() string {
	return m.Loader.View()
}
