package hera

func (model *rootModel) TabByTitle(title string) *commandTab {
	for _, tab := range model.commandTabs {
		if tab.Title != title {
			continue
		}

		return tab
	}

	return nil
}

func (model *rootModel) ActiveTab() *commandTab {
	return model.commandTabs[model.activeTabIndex]
}

func (model *rootModel) NextTab() {
	model.activeTabIndex = (model.activeTabIndex + 1) % len(model.commandTabs)
}

func (model *rootModel) PreviousTab() {
	model.activeTabIndex = (model.activeTabIndex - 1 + len(model.commandTabs)) % len(model.commandTabs)
}

func (model *rootModel) ViewportHeight() int {
	tabHeight := 3
	separatorHeight := 1

	return model.terminalHeight - tabHeight - separatorHeight
}
