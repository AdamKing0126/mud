package main

import (
	"fmt"
	"os"

	"github.com/adamking0126/mud/internal/ui/character"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	fmt.Printf("Character Creation App started. PID: %d\n", os.Getpid())

	model := character.NewCreationModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
