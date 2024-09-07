package login

import (
	"github.com/adamking0126/mud/internal/game/players"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

/*

What are we doing here? How did we get here?
- if the player in the parent component is empty, we know we need to log in the user.  when we are done here, we will set the user object,
which will allow the parent component to switch to a different view

- does the user already have an account or player that they want to log in as?  If so, we will send them through that login flow

- otherwise, send them to the character creation flow

* Let's ignore additional password verification stuff, and take care of that later.
* create a dummy function that retrieves a list of the player objects from the database - right now it will return an empty list
* if the list is empty, we know we need to create a new player character

*/

type componentState int

const (
	stateCharacterPicker componentState = iota
	stateCharacterCreator
)

type Model struct {
	textInput textinput.Model
	Player    *players.Player
	err       error
	state     componentState
}

func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter text"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &Model{
		textInput: ti,
		err:       nil,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// get list of players that belong to the session user? I dunno
	playerCharacters := m.getPlayerCharacters()
	if len(playerCharacters) == 0 {
		m.state = stateCharacterCreator
	} else {
		m.state = stateCharacterPicker
	}

	switch m.state {
	case stateCharacterCreator:
		return m.updateCharacterCreator(msg)
	default:
		return m.updateCharacterPicker(msg)
	}
}

func (m Model) updateCharacterCreator(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

func (m Model) updateCharacterPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// todo implement
	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case stateCharacterCreator:
		return m.viewCharacterCreator()
	default:
		return m.viewCharacterPicker()
	}
}

func (m Model) viewCharacterCreator() string {
	return m.textInput.View()
}

func (m Model) viewCharacterPicker() string {
	// todo implement
	return ""
}

func (m Model) getPlayerCharacters() []*players.Player {
	characters := make([]*players.Player, 0)
	return characters
}
