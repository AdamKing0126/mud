package main

import (
	"fmt"
	"log/slog"
	"os"

  "github.com/charmbracelet/lipgloss"
	"github.com/adamking0126/mud/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type fooModel struct {
	form []components.Component
  selectedElement int
  zoomed bool
	width     int
	height    int
  logger *slog.Logger
  data []*components.FieldData
}

// todo this is a hack that needs to be fixed.
type localFieldData = components.FieldData

func newModel(components []components.Component, logger *slog.Logger) *fooModel {
  data := make([]*localFieldData, len(components))
  model := &fooModel{form: components, zoomed: false, selectedElement: 0, logger: logger, data: data}
  return model
}

func (m *fooModel) Init() tea.Cmd {
  return nil
}

func (m *fooModel) handleZoomUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
    }
  }

  updatedModel, cmd := m.form[m.selectedElement].Update(msg)
  if updatedComponent, ok := updatedModel.(components.Component); ok {
    m.form[m.selectedElement] = updatedComponent
  } else {
    m.logger.Error("Failed to cast updated model to Component")
  }
  return m, cmd
} 

func (m *fooModel) advance(num int, cmd tea.Cmd) (*fooModel, tea.Cmd) {
  m.form[m.selectedElement].SetZoom(false)
  m.form[m.selectedElement].SetHighlighted(false)

  if num < 0 {
    if m.selectedElement == 0 {
      m.selectedElement = len(m.form) - 1
    } else {
      m.selectedElement = m.selectedElement - 1
    }
  } else {
    if m.selectedElement == len(m.form) - 1 {
      m.selectedElement = 0
    } else {
      m.selectedElement = m.selectedElement + 1
    }
  }

  m.form[m.selectedElement].SetHighlighted(true)
  return m, cmd
}

func (m *fooModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
    // TODO figure this out
		// m.component.SetSize(m.width, m.height)
    for _, component := range m.form {
      component.SetSize(m.width, m.height)
    }

    return m, cmd
  case components.SubmitMessage:
    fieldData := msg.Data.(*components.FieldData)
    if fieldData != nil {
      m.data[m.selectedElement] = fieldData
    }
    m.zoomed = false
    m, cmd = m.advance(1, cmd)

	case tea.KeyMsg:
    if m.zoomed {
      return m.handleZoomUpdate(msg)
    }
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
    case "up", "j":
      return m.advance(-1, cmd)
    case "down", "k":
      return m.advance(1, cmd)
    case "enter":
      if !m.zoomed {
        m.zoomed = true
        m.form[m.selectedElement].SetZoom(true)
        return m, nil
      }
		}

    updatedModel, cmd := m.form[m.selectedElement].Update(msg)
    if updatedComponent, ok := updatedModel.(components.Component); ok {
      m.form[m.selectedElement] = updatedComponent
    } else {
      m.logger.Error("Failed to cast updated model to Component")
    }
    return m, cmd

	}

  return m, cmd
}

func (m *fooModel) View() string {
  var zoomedView string
  if m.zoomed {
    zoomedView = m.form[m.selectedElement].ZoomableView()
  }

  if m.form[m.selectedElement].Zoomable() && m.form[m.selectedElement].GetZoom() {
    return zoomedView
  } else {
    var views []string
    for idx, component := range m.form {
      if idx == m.selectedElement {
        if m.zoomed {
          views = append(views, zoomedView)
        } else {
          views = append(views, component.HighlightView())
        }
      } else {
        views = append(views, component.UnfocusedView())
      }
    }
    res := lipgloss.JoinVertical(lipgloss.Top, views...)
    return res
  }
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

  return elementsList
}
