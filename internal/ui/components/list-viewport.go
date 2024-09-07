package components

import (
	"fmt"
	"log/slog"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Item struct {
	title, fieldName, desc, longDesc, id string
}

func (i Item) Title() string           { return i.title }
func (i Item) Description() string     { return i.desc }
func (i Item) LongDescription() string { return i.longDesc }
func (i Item) FilterValue() string     { return i.title }
func (i Item) Id() string              { return i.id }

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
	list     list.Model
	viewport viewport.Model
	state    focusedState
	logger   *slog.Logger
	selected list.Item
	items    []Item
}

func NewListViewportModel(items []Item, logger *slog.Logger) *ListViewportModel {
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
		list:     l,
		viewport: viewport.New(0, 0),
		state:    listFocused,
		logger:   logger,
		items:    items,
	}

	if len(m.list.Items()) > 0 {
		m.updateViewportContent()
	}

	return m
}

func (m *ListViewportModel) SetLogger(logger *slog.Logger) {
	m.logger = logger
}

func (m *ListViewportModel) Init() tea.Cmd {
	return nil
}

func (m *ListViewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right", "h", "l":
			if m.state == listFocused {
				m.state = viewportFocused
				m.updateViewportContent()
				m.logger.Debug("Viewport focused", "msg", msg)
			} else {
				m.state = listFocused
				m.logger.Debug("List focused", "msg", msg)
			}
		}
	}

	if m.state == listFocused {
		oldIndex := m.list.Index()
		m.selected = m.list.SelectedItem()
		m.list, cmd = m.list.Update(msg)

		if m.list.Index() != oldIndex {
			m.logger.Debug("updating Viewport content because list item changed", "msg", msg, "listIndex", m.list.Index(), "oldIndex", oldIndex)
			m.updateViewportContent()
		}
	} else {
		m.logger.Debug("updating Viewport", "msg", msg)
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func (m *ListViewportModel) updateViewportContent() {
	m.logger.Debug("Updating Viewport content", "listIndex", m.list.Index())
	if i, ok := m.list.SelectedItem().(Item); ok {
		content := fmt.Sprintf("Title: %s\n\nDescription: %s", i.Title(), i.LongDescription())

		wrappedContent := wordwrap.String(content, m.viewport.Width)
		m.logger.Debug("updateViewPortContent called", "wrappedContent", wrappedContent)
		m.viewport.SetContent(wrappedContent)
		m.viewport.GotoTop()
	}
}

func (m *ListViewportModel) View() string {
	listView := m.list.View()
	m.logger.Debug("list dimensions", "width", m.list.Width(), "height", m.list.Height())
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
		"id":    selectedItem.id,
    "fieldName": selectedItem.fieldName,
		"title": selectedItem.title,
	}

	return []map[string]string{selectionData}
}
