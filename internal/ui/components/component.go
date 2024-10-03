package components

import (
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "github.com/charmbracelet/bubbles/list"
)


type Component interface {
  list.Item
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	SetSize(width, height int)
	GetSelected() []map[string]string
  GetValue() any
  GoToBeginning(int)
  Zoomable() bool
  SetFocus(bool) 
  HighlightView() string
  UnfocusedView() string
  GetHighlightStyle() lipgloss.Style
  SetZoom(bool) 
  GetZoom() bool
  FilterValue() string
  Title() string
  Description() string
}

type FieldData struct {
  FieldName string
  FieldDescription string
  Value any
}

