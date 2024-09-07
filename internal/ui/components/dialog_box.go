package components

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lipglossList "github.com/charmbracelet/lipgloss/list"
)

var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(2, 2).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginRight(2).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

type SelectedMessage struct {
	Selected []map[string]string
}

type DialogBoxWrapper struct {
	width                int
	height               int
	confirmationQuestion string
	confirmText          string
	cancelText           string
	message              string
	isActive             bool
	component            Component
	logger               *slog.Logger
	acceptButtonActive   bool
}

func NewDialogBox(width, height int, confirmationQuestion string, confirmText string, cancelText string, logger *slog.Logger, component Component, isActive bool) *DialogBoxWrapper {
	return &DialogBoxWrapper{
		width:                width,
		height:               height,
		confirmationQuestion: confirmationQuestion,
		confirmText:          confirmText,
		cancelText:           cancelText,
		component:            component,
		isActive:             isActive,
		logger:               logger,
		acceptButtonActive:   true,
	}

}

func (d *DialogBoxWrapper) Init() tea.Cmd {
	return nil
}

func (d *DialogBoxWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !d.isActive {
		// check for specific key to activate the dialog box
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				d.ToggleActive()
				return d, nil
			}
		}

		// pass the message to the nested component if the dialog box is not activated
		newComponent, cmd := d.component.Update(msg)
		d.component = newComponent.(Component)
		return d, cmd
	}

	// handle dialog box-specific messages when it is active
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "tab", "l", "h":
			d.acceptButtonActive = !d.acceptButtonActive
			return d, nil
		case "enter":
			d.isActive = false
			if d.acceptButtonActive {
				return d, func() tea.Msg {
					return SelectedMessage{Selected: d.component.GetSelected()}
				}
			}
			return d, nil
		case "esc":
			if d.IsActive() {
				d.SetActive(false)
			}
		}
	}

	return d, nil
}

func (d *DialogBoxWrapper) View() string {
	if !d.isActive {
		return d.component.View()
	}

	var okButton, cancelButton string
	if d.acceptButtonActive {
		okButton = activeButtonStyle.Render(d.confirmText)
		cancelButton = buttonStyle.Render(d.cancelText)
	} else {
		okButton = buttonStyle.Render(d.confirmText)
		cancelButton = activeButtonStyle.Render(d.cancelText)
	}

	question := lipgloss.NewStyle().Align(lipgloss.Center).Render(d.confirmationQuestion)
	choiceMessage := lipgloss.NewStyle().Align(lipgloss.Center).Render(d.message)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, choiceMessage, buttons)

	return lipgloss.Place(d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("╬"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}

func (d *DialogBoxWrapper) SetSize(width, height int) {
	d.width = width
	d.height = height
	if d.component != nil {
		d.component.SetSize(d.width, d.height)
	}
}

func (d *DialogBoxWrapper) Activate() {
	d.isActive = true
}

func (d *DialogBoxWrapper) SetActive(isActive bool) {
	d.isActive = isActive
}

func (d *DialogBoxWrapper) ToggleActive() {
	d.isActive = !d.isActive
	if d.isActive {
		d.acceptButtonActive = true
		d.SetMessageFromSelection()
	}
}

func (d *DialogBoxWrapper) IsActive() bool {
	return d.isActive
}

func (d *DialogBoxWrapper) SetMessageFromSelection() {
	selected := d.component.GetSelected()
	if len(selected) == 1 {
		d.message = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(selected[0]["title"])
	} else {
		l := lipglossList.New()
		for _, elem := range selected {
			l.Item(elem["title"])
		}
		d.message = l.String()
	}
}
