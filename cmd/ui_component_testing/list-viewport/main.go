package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/adamking0126/mud/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
  // "github.com/charmbracelet/huh"
)

type model struct {
	component components.Component
	width     int
	height    int
}

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		PaddingLeft(1)
)

func (m model) Init() tea.Cmd {
	return m.component.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.component.SetSize(m.width, m.height)
	}

	var newComponent tea.Model
	newComponent, cmd = m.component.Update(msg)

	if newComponentAsComponent, ok := newComponent.(components.Component); ok {
		m.component = newComponentAsComponent
	} else {
		fmt.Printf("Warning: component.Update() returned a model of type %T that doesn't implement Component\n", newComponent)
	}

	return m, cmd
}

func (m model) View() string {
	title := titleStyle.Render("List and Viewport")
	componentView := m.component.View()

	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		componentView,
	)
	return content + "\n" + statusBarView(m.width)
}

func statusBarView(width int) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("%-*s", width, "q: quit • tab: switch focus"))
}

func main() {
	logFile, err := os.OpenFile("listviewport_testing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))

	items := []components.Item{components.NewItem("1", "character", "Elara the Elven Mage", "Master of arcane arts", `Elara is a centuries-old elven mage, renowned for her mastery of elemental magic and her insatiable thirst for knowledge. With hair like spun silver and eyes that shimmer with arcane power, she cuts an imposing figure in her flowing robes embroidered with mystical sigils.

Born in the ancient elven city of Silvergrove, Elara showed an aptitude for magic from a young age. She spent decades studying under the greatest mages of her people before setting out to explore the wider world, driven by a desire to uncover lost magical secrets and push the boundaries of what was thought possible.

Elara's specialty lies in weaving together different elemental magics to create spectacular and devastating effects. She can call down lightning storms, raise walls of flame, and reshape the very earth beneath her feet. Her most impressive feat to date was single-handedly holding back a tsunami threatening a coastal city, earning her the moniker "The Tide-Turner."

Despite her power, Elara remains humble and driven by an insatiable curiosity. She can often be found poring over ancient tomes in forgotten libraries or exploring dangerous ruins in search of magical artifacts. While generally benevolent, her pursuit of arcane knowledge can sometimes blind her to the consequences of her actions.

Elara serves as both a valuable ally and a potential rival for other magic users, always willing to share her knowledge but also fiercely competitive when it comes to magical prowess. Her ultimate goal is to unravel the fundamental laws of magic itself, a quest that may have far-reaching consequences for the entire world.`),
		components.NewItem("2", "character", "Grimlock the Dwarven Smith", "Legendary craftsman", `Grimlock is a surly but brilliant dwarven smith, capable of forging artifacts of immense power and beauty. His wild, fire-red beard is singed at the edges, and his muscular arms are covered in a tapestry of scars and burn marks – badges of honor from a lifetime at the forge.

Born deep in the heart of the Iron Mountains, Grimlock showed an uncanny talent for metalworking from the moment he could lift a hammer. He apprenticed under the greatest smiths of his clan, quickly surpassing them all. His masterwork – a set of armor that could withstand dragonfire – earned him the title of Master Smith at an unusually young age.

Grimlock's true genius lies in his ability to infuse his creations with magic, a rare talent among dwarves. He combines traditional dwarven craftsmanship with arcane enchantments to create weapons and armor of unparalleled quality. His workshop is always filled with strange alchemical substances, rare metals, and crystals humming with magical energy.

Despite his gruff exterior, Grimlock has a strong sense of honor and takes great pride in his work. He refuses to create for those he deems unworthy, turning away kings and emperors if he feels they lack the character to wield his creations responsibly. For those who do earn his respect, however, he will work tirelessly to craft items perfectly suited to their needs.

Grimlock's ultimate ambition is to forge a weapon or artifact of such power that it will be remembered for millennia to come. He's always on the lookout for rare materials and lost techniques that might help him achieve this goal, making him a valuable ally for adventurers willing to brave dangerous locales in search of such treasures.`),
		components.NewItem("3", "character", "Zephyr the Shapeshifter", "Mysterious changeling", `Zephyr is an enigmatic being, able to assume the form of any creature at will. In their natural state, Zephyr appears as a slender, androgynous figure with iridescent skin that seems to ripple and shift like water. Their eyes are pools of swirling silver, reflecting the countless forms they've worn over their lifetime.

The origins of Zephyr are shrouded in mystery. Some say they are the last of an ancient race of shapeshifters, while others believe they are a magical experiment gone awry. Zephyr themselves either doesn't know or chooses not to reveal the truth of their existence.

Zephyr's shapeshifting abilities go beyond mere physical transformation. They can perfectly mimic the mannerisms, voice, and even surface thoughts of the forms they take. This makes them an unparalleled spy and infiltrator, able to blend seamlessly into any society or situation.

Despite their incredible abilities, Zephyr struggles with questions of identity and belonging. Having lived countless lives in countless forms, they find it difficult to form lasting connections or decide on a true purpose. This inner conflict often manifests as a mercurial personality, with Zephyr's mood and behavior changing as rapidly as their physical form.

Zephyr's motivations are as fluid as their form. Sometimes they use their abilities to help others, taking on the shape of a mighty beast to defend a village or infiltrating a tyrant's court to gather intelligence for rebels. Other times, they might use their powers for personal gain or simply to satisfy their own curiosity about how it feels to live as different creatures.

For those who manage to befriend Zephyr, they can be an invaluable ally, offering unique perspectives and abilities. However, their fluid nature and uncertain loyalties mean that trusting Zephyr completely is always a risk.`),
	}

	listViewportComponent := components.NewListViewportModel(items, logger)

	p := tea.NewProgram(model{component: listViewportComponent}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
