package components

import (
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
  "github.com/charmbracelet/bubbles/list"
  "github.com/google/uuid"
)


type Component interface {
  list.Item

	SetSize(width, height int)
  SetZoom(bool) 
  SetFocus(bool) 
  SetSubmitRecipientId(*uuid.UUID)

  GetId() uuid.UUID
	GetSelected() []map[string]string
  GetValue() any
  GetZoom() bool
  GetHighlightStyle() lipgloss.Style

  // Do we need this?
  GoToBeginning(int)

  Zoomable() bool
  HighlightView() string
  UnfocusedView() string

  /* Bubbletea */
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type FieldData struct {
  FieldName string
  FieldDescription string
  Value any
}

