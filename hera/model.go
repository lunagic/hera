package hera

import (
	"slices"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
	"github.com/lunagic/hera/hera/internal/utils"
)

func newModel(
	userConfig UserConfig,
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
		mouseEnabled: userConfig.EnableMouseByDefault,
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

			slices.SortFunc(commandTabs, func(a *commandTab, b *commandTab) int {
				if a.Title > b.Title {
					return 1
				} else if a.Title < b.Title {
					return -1
				} else {
					return 0
				}
			})

			commandTabs = append(commandTabs, &commandTab{
				triggerRefresh: updateFunc,
				Title:          "help",
				status:         "❔",
				viewport: func() viewport.Model {
					vp := viewport.New(0, 0)
					vp.SetContent(viewHelp())
					return vp
				}(),
			})

			return commandTabs
		}(),
	}
}

type tabPos struct {
	start int // starting x position (column) of a tab header
	width int // width of the rendered tab header
}

type rootModel struct {
	commandTabs    []*commandTab
	mouseEnabled   bool
	activeTabIndex int
	tabPositions   []tabPos
	terminalWidth  int
	terminalHeight int
}

func (model *rootModel) Init() tea.Cmd {
	return tea.Batch(
		func() []tea.Cmd {
			commandTabs := []tea.Cmd{}

			if model.mouseEnabled {
				commandTabs = append(commandTabs, tea.EnableMouseAllMotion)
			} else {
				commandTabs = append(commandTabs, tea.DisableMouse)
			}

			x := 0
			for _, commandTab := range model.commandTabs {
				padding := 7
				commandTabs = append(commandTabs, commandTab.Init())
				width := len(commandTab.Title) + padding
				model.tabPositions = append(model.tabPositions, tabPos{
					start: x,
					width: width,
				})
				x += width
			}

			return commandTabs
		}()...,
	)
}

func (model *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if tea.MouseActionPress == msg.Action {
			for i, pos := range model.tabPositions {
				if msg.Y >= 3 {
					return model, nil
				}
				if msg.X >= pos.start && msg.X < pos.start+pos.width {
					model.activeTabIndex = i
					return model, nil
				}
			}
		}
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
	)
}
