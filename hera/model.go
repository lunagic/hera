package hera

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/lunagic/hera/hera/internal/utils"
)

func newModel(
	config Config,
	updateFunc func(),
	fileChangedFunc func(title string),
) tea.Model {
	config.prime()

	go utils.Watch(
		func(fileName string) error {
			for serviceName, service := range config.Services {
				if !service.shouldTriggerUpdate(fileName) {
					continue
				}

				fileChangedFunc(serviceName)
			}
			return nil
		},
	)

	return &rootModel{
		commandTabs: func() []*commandTab {
			commandTabs := []*commandTab{}
			for serviceName, service := range config.Services {
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
		color.New(color.Reset).Sprint("")+model.ActiveTab().viewport.View()+color.New(color.Reset).Sprint(""),
		model.viewHelp(),
	)
}
