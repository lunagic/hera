package hera

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Italic(true).Faint(true)
	tabStyle  = lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			Bold(true).Border(lipgloss.HiddenBorder())

	activeTabStyle = tabStyle.BorderForeground(lipgloss.NoColor{}).Border(lipgloss.RoundedBorder())
)

func (model *rootModel) viewHelp() string {
	instructions := []string{
		"exit: ctrl-c",
		"restart tab: ctrl-r",
	}

	return helpStyle.Render(strings.Join(instructions, "  ;  "))
}

func (model *rootModel) viewTabs() string {
	tabTitles := []string{}
	for i, tab := range model.commandTabs {
		title := fmt.Sprintf("%s %s", tab.status, tab.Title)
		style := tabStyle
		if i == model.activeTabIndex {
			style = activeTabStyle
		}

		tabTitles = append(tabTitles, style.Render(title))
	}

	tabs := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabTitles...,
	)

	separator := helpStyle.
		Width(model.terminalWidth).
		Height(1).
		Render(
			strings.Repeat("─", model.terminalWidth),
		)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		separator,
	)
}
