package boba

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func NewLoadingModel() LoadingModel {
	return LoadingModel{
		loading: false,
		spinner: spinner.New(),
	}
}

type LoadingModel struct {
	loading bool
	spinner spinner.Model
}

func (loader *LoadingModel) Load() tea.Msg {
	if !loader.loading {
		return loadingMsg{}
	}
	return nil
}

func (loader *LoadingModel) updateLoading(msg tea.Msg, cmds *[]tea.Cmd) spinner.Model {
	switch msg := msg.(type) {
	case loadingMsg:
		loader.loading = true
		*cmds = append(*cmds, loader.spinner.Tick)
	case spinner.TickMsg:
		if loader.loading {
			loader.spinner = UpdateSubmodel(loader.spinner, msg, cmds)
		}
	case MultiSelectorOptionMsg, SelectorOptionsMsg, errMsg:
		loader.loading = false
	}

	return loader.spinner
}
