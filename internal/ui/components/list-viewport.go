package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Item struct {
	title, desc, longDesc string
}

func (i Item) Title() string           { return i.title }
func (i Item) Description() string     { return i.desc }
func (i Item) LongDescription() string { return i.longDesc }
func (i Item) FilterValue() string     { return i.title }

func NewItem(title, desc, longDesc string) Item {
	return Item{title: title, desc: desc, longDesc: longDesc}
}

type focusedState int

const (
	listFocused focusedState = iota
	viewportFocused
)

var (
	focusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("69"))

	unfocusedStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder())
)

type ListViewportModel struct {
	list     list.Model
	viewport viewport.Model
	state    focusedState
}

func NewListViewportModel(items []Item) *ListViewportModel {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.SetShowTitle(false)

	v := viewport.New(0, 0)

	m := &ListViewportModel{
		list:     l,
		viewport: v,
		state:    listFocused,
	}

	if len(items) > 0 {
		m.updateViewportContent()
	}

	return m
}

func (m *ListViewportModel) Init() tea.Cmd {
	return nil
}

func (m *ListViewportModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.state == listFocused {
				m.state = viewportFocused
			} else {
				m.state = listFocused
			}
		}
	}

	if m.state == listFocused {
		oldIndex := m.list.Index()
		m.list, cmd = m.list.Update(msg)

		if m.list.Index() != oldIndex {
			m.updateViewportContent()
		}
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
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

	footer := fmt.Sprintf("[%3.f%%]", m.viewport.ScrollPercent()*100)
	viewportWithFooterView := lipgloss.JoinVertical(lipgloss.Right, viewportView, footer)

	if m.state == listFocused {
		listView = focusedStyle.Render(listView)
		viewportView = unfocusedStyle.Render(viewportView)
	} else {
		listView = unfocusedStyle.Render(listView)
		viewportView = focusedStyle.Render(viewportWithFooterView)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, viewportView)
}

func (m *ListViewportModel) SetSize(width, height int) {
	listWidth := width / 2
	viewportWidth := width - listWidth - 1 // -1 to account for the separator

	m.list.SetSize(listWidth, height)
	m.viewport.Width = viewportWidth
	m.viewport.Height = height - 1 // -1 to account for the footer
}
