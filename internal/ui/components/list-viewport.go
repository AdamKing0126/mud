package components

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
  "github.com/google/uuid"
)

type Item struct {
	title, fieldName, desc, longDesc, id string
}

func (i Item) Title() string           { return i.title }
func (i Item) Description() string     { return i.desc }
func (i Item) LongDescription() string { return i.longDesc }
func (i Item) FilterValue() string     { return i.title }
func (i Item) Id() string              { return i.id }
func (i Item) FieldName() string { return i.fieldName }

func NewItem(id, fieldName, title, desc, longDesc string) Item {
	return Item{title: title, desc: desc, longDesc: longDesc, id: id, fieldName: fieldName}
}

type focusedState int

const (
	listFocused focusedState = iota
	viewportFocused
)

var (
	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69")).
			BorderBottom(true).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true)

	unfocusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
)

type ListViewportModel struct {
  id       uuid.UUID
  submitRecipientId *uuid.UUID
	list     list.Model
	viewport viewport.Model
	state    focusedState
	logger   *slog.Logger
	selected list.Item
	items    []Item
  label    string
  highlightStyle lipgloss.Style
  highlighted bool
  zoomed bool
  dirty bool
}

func NewListViewportModel(submitRecipientId *uuid.UUID, label string, items []Item, highlightStyle lipgloss.Style, logger *slog.Logger) Component {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()

	m := &ListViewportModel{
    id: uuid.New(),
    submitRecipientId: submitRecipientId,
		list:     l,
		viewport: viewport.New(0, 0),
		state:    listFocused,
		logger:   logger,
		items:    items,
    label:    label,
    highlightStyle: highlightStyle,
    highlighted: false,
    zoomed: false,
    dirty: false,
	}

	if len(m.list.Items()) > 0 {
		m.updateViewportContent()
	}

	return m
}

// Bubbletea 
func (m *ListViewportModel) Init() tea.Cmd {
	return nil
}

func (m *ListViewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	m.logger.Debug("ListViewportModel Update", "msg", msg)

  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
      case "left", "right", "h", "l":
        if m.state == listFocused {
          m.state = viewportFocused
          m.updateViewportContent()
        } else {
          m.state = listFocused
        }
      case "ctrl+c", "q":
        return m, tea.Quit 
      case "enter":
        if m.GetZoom() {
          m.dirty = true
          return m, func() tea.Msg {
            return SubmitMessage{
              SenderId: m.id,
              RecipientId: *m.submitRecipientId,
              Data: m.GetValue(),
            }
          }
        } else {
          m.SetZoom(true)
          m.state = listFocused
        }
    }
  }
  
  if m.GetZoom() {
    if m.state == listFocused {
      oldIndex := m.list.Index()
      m.selected = m.list.SelectedItem()
      m.list, cmd = m.list.Update(msg)

      if m.list.Index() != oldIndex {
        m.updateViewportContent()
      }
    } else {
      m.viewport, cmd = m.viewport.Update(msg)
    }
  }

  return m, cmd
}


// Custom

func (m *ListViewportModel) Zoomable() bool {
  return true
}

func (m *ListViewportModel) SetLogger(logger *slog.Logger) {
	m.logger = logger
}

func (m *ListViewportModel) GoToBeginning(_ int) {
  // m.state = listFocused
  m.list.ResetSelected()
}

func (m *ListViewportModel) updateViewportContent() {
	if i, ok := m.list.SelectedItem().(Item); ok {
		content := fmt.Sprintf("Title: %s\n\nDescription: %s", i.Title(), i.LongDescription())

		wrappedContent := wordwrap.String(content, m.viewport.Width)
		m.viewport.SetContent(wrappedContent)
		m.viewport.GotoTop()
	}
}

func (m *ListViewportModel) View() string {
  listView := m.list.View()
  viewportView := m.viewport.View()

  if m.state == listFocused {
    listView = focusedStyle.Render(listView)
    viewportView = unfocusedStyle.Render(viewportView)
  } else {
    listView = unfocusedStyle.Render(listView)
    footer := fmt.Sprintf("[%3.f%%]", m.viewport.ScrollPercent()*100)
    viewportView = focusedStyle.Render(lipgloss.JoinVertical(lipgloss.Right, viewportView, footer))
  }

  return lipgloss.JoinHorizontal(lipgloss.Top, listView, viewportView)
}

func (m *ListViewportModel) SetSize(width, height int) {
	componentHeaderHeight := 3 // 1 for title, 1 for padding

	componentFooterHeight := 1 // for status bar

	componentHeight := height - componentHeaderHeight - componentFooterHeight
	componentWidth := width

	listWidth := componentWidth / 2
	m.list.SetSize(listWidth, componentHeight)
	listComponentWidth := lipgloss.Width(m.list.View())

	viewportWidth := componentWidth - listComponentWidth - 5 // borders

	m.viewport.Width = viewportWidth
	m.viewport.Height = componentHeight - 1 // -1 to account for the "scroll percentage"
}

func (m *ListViewportModel) GetSelected() []map[string]string {
	selectedIndex := m.list.Index()
	selectedItem := m.items[selectedIndex]

	selectionData := map[string]string{
		"id":        selectedItem.id,
		"fieldName": selectedItem.fieldName,
		"title":     selectedItem.title,
	}

	return []map[string]string{selectionData}
}

func (m *ListViewportModel) GetValue() any {
  if i, ok := m.list.SelectedItem().(Item); ok {
    return &FieldData{
      FieldName: i.FieldName(),
      Value: i.Id(),
      FieldDescription: i.Title(),
    }
  }
  return nil
}

func (m *ListViewportModel) SetFocus(bool) {
}

func (m *ListViewportModel) GetHighlightStyle() lipgloss.Style {
  return m.highlightStyle
}

func (m *ListViewportModel) HighlightView() string {
  var ret string

  if m.dirty {
    value := m.GetValue()
    fieldData := value.(*FieldData)
    ret = fmt.Sprintf("%s (enter to change)", fieldData.FieldDescription)
  } else {
    ret = "(enter to set)"
  }

  highlightText := m.GetHighlightStyle().Render(ret)
  labelText := fmt.Sprintf("%s: ", m.label)
  return lipgloss.JoinHorizontal(lipgloss.Top, labelText, highlightText)
}

func (m *ListViewportModel) UnfocusedView() string {
  if m.dirty {
    val := m.GetValue().(*FieldData)
    return fmt.Sprintf("%s: %s", m.label, val.FieldDescription)
  }
  return fmt.Sprintf("%s: (empty)", m.label)
}

func (m *ListViewportModel) SetZoom(zoomed bool) {
  m.zoomed = zoomed
}

func (m *ListViewportModel) GetZoom() bool {
  return m.zoomed
}

func (m *ListViewportModel) GetId() uuid.UUID {
  return m.id
}

func (m *ListViewportModel) SetSubmitRecipientId(submitRecipientId *uuid.UUID) {
  m.submitRecipientId = submitRecipientId
}

/* statisfy list.Item interface */
func (m *ListViewportModel) FilterValue() string {
  return ""
}

func (m *ListViewportModel) Title() string {
  return ""
}

func (m *ListViewportModel) Description() string {
  return ""
}
