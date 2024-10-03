package components

import (
  "fmt"
  "log/slog"
  "github.com/charmbracelet/bubbles/textinput"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

type InputComponent struct {
  textInput textinput.Model
  label string
  logger *slog.Logger
  width int
  height int
  highlightStyle lipgloss.Style
  highlighted bool
  zoomed bool
}

func NewInputComponent(highlightStyle lipgloss.Style, label string, placeholder string, charlimit int, width int, logger *slog.Logger) *InputComponent {
  ti := textinput.New()
  ti.Placeholder = placeholder
  ti.CharLimit = charlimit
  ti.Width = width

  return &InputComponent{
    textInput: ti,
    logger: logger,
    highlightStyle: highlightStyle,
    highlighted: false,
    label: label,
    zoomed: false,
  }
}

// These functions satisfy the Bubbletea interface
func (i *InputComponent) Init() tea.Cmd {
  return textinput.Blink
}

func (i *InputComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  i.logger.Debug("InputComponent received msg", "msg", msg)
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "enter":
      if i.GetZoom() {
        value := i.GetValue()
        i.SetZoom(false)
        i.textInput.Blur()
        return i, func() tea.Msg {
          fieldData := &FieldData{
            FieldName: "foo",
            Value: value,
            FieldDescription: "bar",
          }
          return SubmitMessage{Data: fieldData}
          }
      } else {
        i.SetZoom(true)
      }
    } 
  }

  if i.GetZoom() {
    i.textInput, cmd = i.textInput.Update(msg)
  }

  i.logger.Debug(fmt.Sprintf("InputComponent value: %s", i.textInput.Value()))
  return i, cmd
}

func (i *InputComponent) View() string {
  return fmt.Sprintf("%s ", i.label) + i.textInput.View()
}

// Component Interface-specific function 
func (i *InputComponent) SetFocus(focus bool) {
  if focus {
    i.textInput.Focus()
  } else { 
    i.textInput.Blur()
  }
}

func (i *InputComponent) GetSelected() []map[string]string {
  return nil
}

func (i *InputComponent) GoToBeginning(_ int) {
}

func (i *InputComponent) SetSize(width, height int) {
  i.height = height
  i.width = width
}

func (i *InputComponent) GetValue() any {
  return i.textInput.Value()
}

func (i *InputComponent) Zoomable() bool {
  return false
}

func (i *InputComponent) GetHighlightStyle() lipgloss.Style {
  return i.highlightStyle
}

func (i *InputComponent) HighlightView() string {
  val := i.textInput.Value()
  var ret string
  if val == "" {
    ret = "(enter) to set"
  } else {
    ret = fmt.Sprintf("%s (enter to change)", val)
  }
  highlightText := i.GetHighlightStyle().Render(ret)

  return lipgloss.JoinHorizontal(lipgloss.Top, fmt.Sprintf("%s: ", i.label), highlightText)
}

func (i *InputComponent) UnfocusedView() string {
  val := i.textInput.Value()
  ret := fmt.Sprintf("%s: ", i.label)
  if val == "" {
    ret += "(empty)"
  }  else {
    ret += val
  }
  return ret
}

func (i *InputComponent) SetZoom(zoomed bool) {
  i.zoomed = zoomed
  i.highlighted = !zoomed
  if zoomed {
    i.textInput.Focus()
  } else {
    i.textInput.Blur()
  }
}

func (i *InputComponent) GetZoom() bool {
  return i.zoomed
}

func (i *InputComponent) FilterValue() string {
  return i.textInput.View()
}

func (i *InputComponent) Title() string {
  return ""
}

func (i *InputComponent) Description() string {
  return ""
}

