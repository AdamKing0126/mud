package character

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreationModel struct {
	nameInput    textinput.Model
	classInput   textinput.Model
	raceInput    textinput.Model
	currentField int
	err          error
}

func NewCreationModel() *CreationModel {
	m := &CreationModel{
		nameInput:  textinput.New(),
		classInput: textinput.New(),
		raceInput:  textinput.New(),
	}

	m.nameInput.Placeholder = "Enter character name"
	m.classInput.Placeholder = "Enter character class"
	m.raceInput.Placeholder = "Enter character race"

	m.nameInput.Focus()

	return m
}

func (m *CreationModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *CreationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			if m.currentField == 2 {
				// Character creation complete
				return m, tea.Quit
			}
			m.currentField++
			if m.currentField > 2 {
				m.currentField = 2
			}
		case "tab":
			m.currentField = (m.currentField + 1) % 3
		case "shift+tab":
			m.currentField = (m.currentField - 1 + 3) % 3
		default:
			// Add the input to the current field
			switch m.currentField {
			case 0:
				m.nameInput.SetValue(m.nameInput.Value() + msg.String())
			case 1:
				m.classInput.SetValue(m.classInput.Value() + msg.String())
			case 2:
				m.raceInput.SetValue(m.raceInput.Value() + msg.String())
			}
		}
	}

	return m, nil
}

func (m *CreationModel) View() string {
	return fmt.Sprintf(
		"Character Creation\n\n"+
			"Name: %s\n"+
			"Class: %s\n"+
			"Race: %s\n\n"+
			"Current field: %d\n"+
			"Press Tab to switch fields, Enter to confirm, q to quit",
		m.nameInput.Value(),
		m.classInput.Value(),
		m.raceInput.Value(),
		m.currentField,
	)
}
