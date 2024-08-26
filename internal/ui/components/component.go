package components

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	SetSize(width, height int)
}
