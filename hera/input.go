package hera

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (model *rootModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab", "right":
		model.NextTab()
	case "shift+tab", "left":
		model.PreviousTab()
	case "up":
		model.ActiveTab().viewport.LineUp(1)
	case "down":
		model.ActiveTab().viewport.LineDown(1)
	case "pgup":
		model.ActiveTab().viewport.LineUp(model.ViewportHeight())
	case "pgdown":
		model.ActiveTab().viewport.LineDown(model.ViewportHeight())
	case "ctrl+l":
		model.ActiveTab().commandOutput = ""
		model.ActiveTab().viewport.SetContent("")
	case "ctrl+b":
		model.ActiveTab().viewport.GotoBottom()
	case "ctrl+r":
		return model, model.ActiveTab().Init()
	case "ctrl+a":
		if model.mouseEnabled {
			return model, tea.DisableMouse
		}

		model.mouseEnabled = true

		return model, tea.EnableMouseAllMotion
	case "ctrl+c", "q":
		for _, tab := range model.commandTabs {
			if tab.processTracker == nil {
				continue
			}
			tab.processTracker.KillAll()
		}

		return model, tea.Quit
	}

	return model, nil
}
