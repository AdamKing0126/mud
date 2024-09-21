package components

import (
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)


type Component interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
  ZoomableView() string
	SetSize(width, height int)
	GetSelected() []map[string]string
  GetValue() any
  GoToBeginning(int)
  Zoomable() bool
  SetFocus(bool) 
  SetHighlighted(bool)
  HighlightView() string
  UnfocusedView() string
  GetHighlightStyle() lipgloss.Style
  SetZoom(bool) 
  GetZoom() bool
}

type FieldData struct {
  FieldName string
  FieldDescription string
  Value any
}

