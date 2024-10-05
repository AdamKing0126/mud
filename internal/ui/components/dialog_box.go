package components

import (
  "fmt"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lipglossList "github.com/charmbracelet/lipgloss/list"
  "github.com/google/uuid"
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

type DialogBoxWrapper struct {
  id                   uuid.UUID
  submitRecipientId    *uuid.UUID
	Width                int
	Height               int
	ConfirmationQuestion string
	ConfirmText          string
	CancelText           string
	Message              string
	IsActive             bool
	Component            Component
	Logger               *slog.Logger
	AcceptButtonActive   bool
  Data                 []FieldData
}

func NewDialogBox(submitRecipientId *uuid.UUID, Width, Height int, ConfirmationQuestion string, ConfirmText string, CancelText string, logger *slog.Logger, component Component, IsActive bool) *DialogBoxWrapper {
	return &DialogBoxWrapper{
    id:                   uuid.New(),
    submitRecipientId:    submitRecipientId,
		Width:                Width,
		Height:               Height,
		ConfirmationQuestion: ConfirmationQuestion,
		ConfirmText:          ConfirmText,
		CancelText:           CancelText,
		Component:            component,
		IsActive:             IsActive,
		Logger:               logger,
		AcceptButtonActive:   true,
	}

}

func (d *DialogBoxWrapper) Init() tea.Cmd {
	return nil
}

func (d *DialogBoxWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.Logger.Debug("DialogBoxWrapper Update", "msg", fmt.Sprintf("%T", msg))

	if !d.IsActive {
    d.Logger.Debug("DialogBoxWrapper message received:", "msg", fmt.Sprintf("%T", msg)) 
		switch msg := msg.(type) {
		case SubmitMessage:
			d.Logger.Debug("DialogBoxWrapper received SubmitMessage when Inactive", "msg", msg)
      d.IsActive = true
      data, ok := msg.Data.([]FieldData)
      if !ok {
        d.Logger.Error("Type assertion failed: msg.Data is not of type []FieldData")
        return d, nil
      }
      d.Data = data
      return d, nil
		}

		newComponent, cmd := d.Component.Update(msg)
		d.Component = newComponent.(Component)

		d.Logger.Debug("DialogBoxWrapper Returning", "model", d)
		return d, cmd
	}

	// handle dialog box-specific messages when it is active
	switch msg := msg.(type) {
	case tea.KeyMsg:
		d.Logger.Debug("DialogBoxWrapper Received KeyMessage when Active", "msg", msg)
		switch msg.String() {
		case "left", "right", "tab", "l", "h":
			d.AcceptButtonActive = !d.AcceptButtonActive
			d.Logger.Debug("returning", "ret", "nothing2")
			return d, nil
		case "enter":
      d.Logger.Debug("DialogBoxWrapper.AcceptButtonActive", "val", d.AcceptButtonActive)
			if d.AcceptButtonActive {
        d.Logger.Debug("DialogBoxWrapper handling 'submit'", "msg", msg)
        d.IsActive = false
        // TODO: return a DialogBoxSubmitMessage or something
        return d, tea.Quit
			}
      d.AcceptButtonActive = !d.AcceptButtonActive
      d.Component.GoToBeginning(0)
			d.IsActive = false
			return d, nil
		case "esc":
			if d.IsActive {
				d.SetActive(false)
			}
		}
	}

	d.Logger.Debug("returning", "ret", "nothing")
	return d, nil
}

func (d *DialogBoxWrapper) View() string {
	if !d.IsActive {
		return d.Component.View()
	}

	var okButton, cancelButton string
	if d.AcceptButtonActive {
		okButton = activeButtonStyle.Render(d.ConfirmText)
		cancelButton = buttonStyle.Render(d.CancelText)
	} else {
		okButton = buttonStyle.Render(d.ConfirmText)
		cancelButton = activeButtonStyle.Render(d.CancelText)
	}

	question := lipgloss.NewStyle().Align(lipgloss.Center).Render(d.ConfirmationQuestion)
  var messages []string
  for _, entry := range d.Data {
    messages = append(messages, fmt.Sprintf("* %s: %s (id %s)", entry.FieldName, entry.FieldDescription, entry.Value))
  }
	choiceMessage := lipgloss.NewStyle().Align(lipgloss.Center).Render(lipgloss.JoinVertical(lipgloss.Left, messages...))
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, choiceMessage, buttons)

	return lipgloss.Place(d.Width, d.Height,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("╬"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}

func (d *DialogBoxWrapper) SetSize(Width, Height int) {
	d.Width = Width
	d.Height = Height
	if d.Component != nil {
		d.Component.SetSize(d.Width, d.Height)
	}
}

func (d *DialogBoxWrapper) Activate() {
	d.IsActive = true
}

func (d *DialogBoxWrapper) SetActive(IsActive bool) {
	d.IsActive = IsActive
}

func (d *DialogBoxWrapper) ToggleActive() {
	d.IsActive = !d.IsActive
	if d.IsActive {
		d.AcceptButtonActive = true
		d.SetMessageFromSelection()
	}
}

func (d *DialogBoxWrapper) SetMessageFromSelection() {
	selected := d.Component.GetSelected()
	if len(selected) == 1 {
		d.Message = lipgloss.NewStyle().Width(50).Align(lipgloss.Left).Render(selected[0]["title"])
	} else {
		l := lipglossList.New()
		for _, elem := range selected {
			l.Item(elem["title"])
		}
		d.Message = l.String()
	}
}

func (m *DialogBoxWrapper) GetValue() any {
  return m.Data
}

func (m *DialogBoxWrapper) GetId() uuid.UUID {
  return m.id
}
