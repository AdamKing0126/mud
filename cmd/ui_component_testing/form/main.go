package main

import (
	"fmt"
	"log/slog"
	"os"
  "io"

  "github.com/charmbracelet/lipgloss"
	"github.com/adamking0126/mud/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/list"
)

// todo this is a hack that needs to be fixed.
type localFieldData = components.FieldData

func updateFunc(msg tea.Msg, m *list.Model) tea.Cmd {
    index := m.Index()
    component := m.SelectedItem().(components.Component)

    if component.Zoomable() {
      component.SetSize(m.Width(), m.Height())
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.String() {
      case "enter":
        if !component.GetZoom() && component.Zoomable(){
          height := m.Height()
          delegate := NewItemDelegate(height)
          m.SetDelegate(delegate)
          component.SetZoom(true)
        }
      case "up", "down", "j", "k":
        if component.GetZoom() {
          _, cmd := component.Update(msg)
          return cmd
        }
      }
    }

    if component.GetZoom() {
      updatedComponent, newMsg := component.Update(msg)
      listItemComponent, ok := updatedComponent.(components.Component)
      if !ok {
        return nil
      }

      m.SetItem(index, listItemComponent)
      return newMsg
    }
  return nil
}

type itemDelegate struct {
  UpdateFunc func(mgs tea.Msg, m *list.Model) tea.Cmd
  height int
}

func (d itemDelegate) Render (w io.Writer, m list.Model, index int, item list.Item) {

  component := item.(components.Component)

  isSelected := index == m.Index()

  if !isSelected {
    fmt.Fprint(w, component.UnfocusedView())
  } else {
    if component.GetZoom() {
      fmt.Fprint(w, component.View())
    } else {
      fmt.Fprint(w, component.HighlightView())
    }
  }
}

func (d itemDelegate) Height() int {
  return d.height
}

func (d itemDelegate) Spacing() int {
  return 1
}

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.UpdateFunc(msg, m)
}

func NewItemDelegate(height int) itemDelegate {
  delegate := itemDelegate{height: height}
  delegate.UpdateFunc = updateFunc
  return delegate
}

type FormModel struct {
  list list.Model
  zoomed bool
	width     int
	height    int
  logger *slog.Logger
  data []*components.FieldData
}

func newModel(components []components.Component, logger *slog.Logger) *FormModel {
  items := make([]list.Item, len(components))
  for i, component := range components {
    items[i] = component
      
  }

  itemDelegate := NewItemDelegate(1)
  l := list.New(items, itemDelegate, 0, 0)

  l.Title = "Choose Your Fighter"

  data := make([]*localFieldData, len(components))
  model := &FormModel{list: l, zoomed: false, logger: logger, data: data}
  return model
}

func (m *FormModel) SetZoom(zoomed bool) {
  m.zoomed = zoomed
}

func (m *FormModel) GetZoom() bool {
  return m.zoomed
}

func (m *FormModel) Init() tea.Cmd {
  return nil
}

func (m *FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
    m.list.SetSize(m.width, m.height-5) // TODO fix this -5 stuff ! 
  case components.SubmitMessage:
    fieldData := msg.Data.(*components.FieldData)
    if fieldData != nil {
      m.data[m.list.Index()] = fieldData
    }
    m.SetZoom(false)
    component := m.list.SelectedItem().(components.Component)
    component.SetZoom(false)
    itemDelegate := NewItemDelegate(1)
    m.list.SetDelegate(itemDelegate)
    m.list.CursorDown()
	case tea.KeyMsg:
    if m.GetZoom() {
      component := m.list.SelectedItem()
      if updatedComponent, ok := component.(components.Component); ok {
          component, cmd := updatedComponent.Update(msg)
          listComponent := component.(list.Item)
          m.list.SetItem(m.list.Index(), listComponent)
          return m, cmd
      }
      return m, cmd
    } else {
      switch msg.String() {
      case "ctrl+c", "q":
        return m, tea.Quit
      case "enter":
        m.SetZoom(true)
        selected := m.list.SelectedItem().(components.Component)
        selected.SetZoom(true)

        if selected.Zoomable() {
          delegate := NewItemDelegate(m.height)
          m.list.SetDelegate(delegate)
          m.list.SetItem(m.list.Index(), selected)
        }
        return m,cmd
      }
    }
	}

  m.list, cmd = m.list.Update(msg)
  return m, cmd
}

func (m *FormModel) View() string {
  return m.list.View()
}

func main() {
	logFile, err := os.OpenFile("form_testing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))

	elementsList := CreateFormFields(logger)
  m := newModel(elementsList, logger) 

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}

}

func CreateFormFields(logger *slog.Logger) []components.Component {
  highlightStyle := lipgloss.NewStyle().Bold(true).Foreground(
  lipgloss.Color("15")).Background(lipgloss.Color("22"))

  playerNameInput := components.NewInputComponent(highlightStyle, "Name", "What's your name?", 255, 30, logger)
  fieldName := "Character"


  items := []components.Item{components.NewItem("1", fieldName, "Elara the Elven Mage", "Master of arcane arts", `Elara is a centuries-old elven mage, renowned for her mastery of elemental magic and her insatiable thirst for knowledge. With hair like spun silver and eyes that shimmer with arcane power, she cuts an imposing figure in her flowing robes embroidered with mystical sigils.

Born in the ancient elven city of Silvergrove, Elara showed an aptitude for magic from a young age. She spent decades studying under the greatest mages of her people before setting out to explore the wider world, driven by a desire to uncover lost magical secrets and push the boundaries of what was thought possible.

Elara's specialty lies in weaving together different elemental magics to create spectacular and devastating effects. She can call down lightning storms, raise walls of flame, and reshape the very earth beneath her feet. Her most impressive feat to date was single-handedly holding back a tsunami threatening a coastal city, earning her the moniker "The Tide-Turner."

Despite her power, Elara remains humble and driven by an insatiable curiosity. She can often be found poring over ancient tomes in forgotten libraries or exploring dangerous ruins in search of magical artifacts. While generally benevolent, her pursuit of arcane knowledge can sometimes blind her to the consequences of her actions.

Elara serves as both a valuable ally and a potential rival for other magic users, always willing to share her knowledge but also fiercely competitive when it comes to magical prowess. Her ultimate goal is to unravel the fundamental laws of magic itself, a quest that may have far-reaching consequences for the entire world.`),
		components.NewItem("2", fieldName, "Grimlock the Dwarven Smith", "Legendary craftsman", `Grimlock is a surly but brilliant dwarven smith, capable of forging artifacts of immense power and beauty. His wild, fire-red beard is singed at the edges, and his muscular arms are covered in a tapestry of scars and burn marks – badges of honor from a lifetime at the forge.

Born deep in the heart of the Iron Mountains, Grimlock showed an uncanny talent for metalworking from the moment he could lift a hammer. He apprenticed under the greatest smiths of his clan, quickly surpassing them all. His masterwork – a set of armor that could withstand dragonfire – earned him the title of Master Smith at an unusually young age.

Grimlock's true genius lies in his ability to infuse his creations with magic, a rare talent among dwarves. He combines traditional dwarven craftsmanship with arcane enchantments to create weapons and armor of unparalleled quality. His workshop is always filled with strange alchemical substances, rare metals, and crystals humming with magical energy.

Despite his gruff exterior, Grimlock has a strong sense of honor and takes great pride in his work. He refuses to create for those he deems unworthy, turning away kings and emperors if he feels they lack the character to wield his creations responsibly. For those who do earn his respect, however, he will work tirelessly to craft items perfectly suited to their needs.

Grimlock's ultimate ambition is to forge a weapon or artifact of such power that it will be remembered for millennia to come. He's always on the lookout for rare materials and lost techniques that might help him achieve this goal, making him a valuable ally for adventurers willing to brave dangerous locales in search of such treasures.`),
		components.NewItem("3", fieldName, "Zephyr the Shapeshifter", "Mysterious changeling", `Zephyr is an enigmatic being, able to assume the form of any creature at will. In their natural state, Zephyr appears as a slender, androgynous figure with iridescent skin that seems to ripple and shift like water. Their eyes are pools of swirling silver, reflecting the countless forms they've worn over their lifetime.

The origins of Zephyr are shrouded in mystery. Some say they are the last of an ancient race of shapeshifters, while others believe they are a magical experiment gone awry. Zephyr themselves either doesn't know or chooses not to reveal the truth of their existence.

Zephyr's shapeshifting abilities go beyond mere physical transformation. They can perfectly mimic the mannerisms, voice, and even surface thoughts of the forms they take. This makes them an unparalleled spy and infiltrator, able to blend seamlessly into any society or situation.

Despite their incredible abilities, Zephyr struggles with questions of identity and belonging. Having lived countless lives in countless forms, they find it difficult to form lasting connections or decide on a true purpose. This inner conflict often manifests as a mercurial personality, with Zephyr's mood and behavior changing as rapidly as their physical form.

Zephyr's motivations are as fluid as their form. Sometimes they use their abilities to help others, taking on the shape of a mighty beast to defend a village or infiltrating a tyrant's court to gather intelligence for rebels. Other times, they might use their powers for personal gain or simply to satisfy their own curiosity about how it feels to live as different creatures.

For those who manage to befriend Zephyr, they can be an invaluable ally, offering unique perspectives and abilities. However, their fluid nature and uncertain loyalties mean that trusting Zephyr completely is always a risk.`),
	}

	listViewportComponent := components.NewListViewportModel(fieldName, items, highlightStyle, logger)
  elementsList := []components.Component{playerNameInput, listViewportComponent}

	// _ = components.NewListViewportModel(fieldName, items, highlightStyle, logger)
  // elementsList := []components.Component{playerNameInput}

  return elementsList
}

