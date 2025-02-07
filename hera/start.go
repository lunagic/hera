package hera

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunagic/hera/hera/internal/utils"
)

func Start() {
	config := Config{}

	if err := utils.MustReadYamlConfig(
		[]string{
			".config/hera.yaml",
			"hera.yaml",
		},
		&config,
	); err != nil {
		log.Fatal(err)
	}

	program := tea.NewProgram(nil)
	model := newModel(
		config,
		func() {
			program.Send(eventCommandOutput{})
		},
		func(title string) {
			program.Send(eventFileChanged{ServiceName: title})
		},
	)
	program = tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)
	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
