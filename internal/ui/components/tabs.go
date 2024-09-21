package components

import (
	"fmt"
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tabStyle = lipgloss.NewStyle().
			Border(tabBorder, true).
			BorderForeground(highlightColor).
			Padding(0, 1)

	activeTabBorder      = tabBorder
	firstTabBorder       = tabBorder
	activeFirstTabBorder = tabBorder

	tabGapBorder = lipgloss.Border{
		Bottom:      "─",
		Right:       "",
		BottomRight: "┐",
	}

	tab            lipgloss.Style
	activeTab      lipgloss.Style
	firstTab       lipgloss.Style
	activeFirstTab lipgloss.Style
	tabGap         lipgloss.Style
)

func statusBarView(width int) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("%-*s", width, "<enter>: select • q: quit • ↑/k up • ↓/j down ←/h, →/l: switch focus, <tab>: switch tabs"))
}

func init() {
	activeTabBorder.Bottom = " "
	activeTabBorder.BottomLeft = "┘"
	activeTabBorder.BottomRight = "└"

	firstTabBorder.BottomLeft = "├"

	activeFirstTabBorder.Bottom = ""
	activeFirstTabBorder.BottomLeft = "│"
	activeFirstTabBorder.BottomRight = "└"

	firstTab = tabStyle.Border(firstTabBorder, true)
	activeFirstTab = tabStyle.Border(activeFirstTabBorder, true)

	tab = tabStyle
	activeTab = tabStyle.Border(activeTabBorder, true)

	tabGap = tabStyle.Border(tabGapBorder, true).
		BorderTop(false).
		BorderLeft(false)
}

func NewTabsModel(tabs []string, tabContent []Component, logger *slog.Logger) *TabsModel {
	return &TabsModel{Tabs: tabs, TabContent: tabContent, logger: logger}
}

type TabsModel struct {
	Tabs       []string
	TabContent []Component
	activeTab  int
	width      int
	height     int
	logger     *slog.Logger
}

func (m *TabsModel) SetSize(width, height int) {
	tabsRowHeight := m.getTabsRowHeight()

	borderWidth, borderHeight := m.getBorderSize()

	m.width = width - 2*borderWidth
	m.height = height - 2*borderHeight

	contentHeight := m.height - tabsRowHeight
	contentWidth := m.width - 1 // TODO figure out why I have to do this

	for _, component := range m.TabContent {
		component.SetSize(contentWidth, contentHeight)
	}
}

func (m TabsModel) Init() tea.Cmd {
	return nil
}

func (m *TabsModel) GoToBeginning(idx int) {
  if idx >= len(m.Tabs) {
    m.activeTab = 0
  } else {
    m.activeTab = idx
  }
}

func (m *TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logger.Debug("TabsModel Update", "msg", msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.activeTab == len(m.Tabs)-1 {
				m.activeTab = 0
			} else {
				m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			}
      return m, nil
    case "shift+tab":
      if m.activeTab == 0 {
        m.activeTab = len(m.Tabs) - 1
      } else {
        m.activeTab = m.activeTab - 1
      }
      return m, nil
		case "enter":
			if m.activeTab == len(m.Tabs)-1 {
				value := m.GetValue()
				m.logger.Debug("TabsModel Submitting message to DialogBoxWrapper(hopefully)", "value", value)
				return m, tea.Cmd(func() tea.Msg {
					return SubmitMessage{Data: value}
				})
			}
			m.activeTab = m.activeTab + 1
			return m, nil
		}
	}

	var cmd tea.Cmd
	updatedModel, cmd := m.TabContent[m.activeTab].Update(msg)

	if updatedCompoent, ok := updatedModel.(Component); ok {
		m.TabContent[m.activeTab] = updatedCompoent
	} else {
		fmt.Printf("Warning: component.Update() returned a model of type %T that doesn't implement Component\n", updatedModel)
	}

	return m, cmd
}

func (m *TabsModel) buildTabsRow() string {
	var renderedTabs []string
	for i, t := range m.Tabs {
		if i == 0 {
			if m.activeTab == i {
				renderedTabs = append(renderedTabs, activeFirstTab.Render(t))
			} else {
				renderedTabs = append(renderedTabs, firstTab.Render(t))
			}
		} else {
			if m.activeTab == i {
				renderedTabs = append(renderedTabs, activeTab.Render(t))
			} else {
				renderedTabs = append(renderedTabs, tab.Render(t))
			}
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// todo: why do I have to subtract 1?
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row)-1)))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
}

func (m *TabsModel) buildBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.NormalBorder())
}

func (m *TabsModel) getBorderSize() (int, int) {
	borderStyle := m.buildBorderStyle()
	borderWidth := borderStyle.GetHorizontalFrameSize()
	borderHeight := borderStyle.GetVerticalFrameSize()
	return borderWidth, borderHeight
}

func (m *TabsModel) getTabsRowHeight() int {
	tabsRow := m.buildTabsRow()
	return lipgloss.Height(tabsRow)
}

func (m *TabsModel) View() string {
	tabsRow := m.buildTabsRow()
	// calculate remaining height for the content
	contentHeight := m.height - lipgloss.Height(tabsRow)

	baseBorderStyle := m.buildBorderStyle()

	// render content
	contentStyle := baseBorderStyle.
		Width(m.width).
		Height(contentHeight).
		BorderTop(false).
		BorderForeground(lipgloss.Color("69")).
		Padding(0, 1)

	content := contentStyle.Render(m.TabContent[m.activeTab].View())
	statusBar := statusBarView(m.width)
	result := lipgloss.JoinVertical(lipgloss.Left, tabsRow, content, statusBar)
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *TabsModel) GetSelected() []map[string]string {
	var data []map[string]string
	for _, component := range m.TabContent {
		data = append(data, component.GetSelected()...)
	}
	return data
}

func (m *TabsModel) GetValue() any {
	var valueList []FieldData
	for _, elem := range m.TabContent {
    val, ok := elem.GetValue().(FieldData)
    if !ok {
      m.logger.Error("Type assertion failed: msg.Data is not of type FieldData")
    }
		valueList = append(valueList, val)
	}
  m.logger.Debug("TabsModel returning data", "value", valueList)
	return valueList
}

