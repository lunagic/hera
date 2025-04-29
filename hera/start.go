package hera

import (
	"fmt"
	"log"
	"os"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lunagic/hera/hera/internal/utils"
)

func Start(args ...string) {
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

	// If services names were provided, filter down to just those
	if len(args) > 0 {
		for name := range config.Services {
			if slices.Contains(args, name) {
				continue
			}

			delete(config.Services, name)
		}
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
