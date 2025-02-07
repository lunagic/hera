package hera

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lunagic/hera/hera/internal/utils"
)

func newModel(
	config Config,
	updateFunc func(),
	fileChangedFunc func(title string),
) tea.Model {
	return &rootModel{
		commandTabs: func() []*commandTab {
			commandTabs := []*commandTab{}
			for serviceName, service := range config.Services {
				go utils.Watch(
					service.Watch,
					service.Exclude,
					func(fileName string) error {
						fileChangedFunc(serviceName)
						return nil
					},
				)
				commandTabs = append(
					commandTabs,
					newCommandTab(
						serviceName,
						service.Command,
						updateFunc,
					),
				)
			}

			return commandTabs
		}(),
	}
}

type rootModel struct {
	commandTabs    []*commandTab
	activeTabIndex int
	terminalWidth  int
	terminalHeight int
}

func (model *rootModel) Init() tea.Cmd {
	return tea.Batch(
		func() []tea.Cmd {
			commandTabs := []tea.Cmd{}
			for _, commandTab := range model.commandTabs {
				commandTabs = append(commandTabs, commandTab.Init())
			}

			return commandTabs
		}()...,
	)
}

func (model *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return model.handleKeyMsg(msg)
	case eventFileChanged:
		return model, model.TabByTitle(msg.ServiceName).Init()
	case tea.WindowSizeMsg:
		model.terminalWidth = msg.Width
		model.terminalHeight = msg.Height

		for _, tab := range model.commandTabs {
			tab.viewport.Width = model.terminalWidth
			tab.viewport.Height = model.ViewportHeight()
		}
	}

	return model, nil
}

func (model *rootModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		model.viewTabs(),
		model.ActiveTab().viewport.View(),
		model.viewHelp(),
	)
}
