package main

import (
	"fmt"
	"os"

	"github.com/adamking0126/mud/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 2 // 1 for title, 1 for padding
		footerHeight := 1 // for status bar
		m.component.SetSize(m.width, m.height-headerHeight-footerHeight)
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
	content := lipgloss.JoinVertical(lipgloss.Left,
		title,
		m.component.View(),
	)
	return content + "\n" + statusBarView(m.width)
}

func statusBarView(width int) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("%-*s", width, "q: quit • tab: switch focus"))
}

func main() {
	itemData := []map[string]string{
		{
			"title":    "The Enchanted Forest",
			"desc":     "A mystical woodland realm",
			"longDesc": "The Enchanted Forest is a realm of wonder and magic, where ancient trees whisper secrets to those who listen. Sunlight filters through a canopy of emerald leaves, casting dappled shadows on a carpet of moss and wildflowers. Ethereal mists curl around gnarled roots, concealing hidden pathways and forgotten ruins.\n\nMysterious creatures flit between the trees – iridescent fairies, mischievous sprites, and wise old talking animals. The air is thick with the scent of pine and wild herbs, and the sound of tinkling laughter echoes from unseen sources.\n\nAt the heart of the forest stands the ancient Tree of Life, its branches reaching towards the heavens and its roots delving deep into the earth. Legend says that those who touch its bark can glimpse visions of past and future.\n\nBrave adventurers who enter the Enchanted Forest may find themselves facing magical trials, solving riddles posed by cryptic beings, or stumbling upon hidden groves of power. But beware – not all that glitters is benevolent, and dark forces lurk in the deepest shadows of this mystical realm.",
		},
		{
			"title":    "The Sunken City",
			"desc":     "An underwater metropolis",
			"longDesc": "Beneath the rolling waves of the ocean lies the Sunken City, a once-great metropolis now claimed by the sea. Massive structures of coral-encrusted marble and weathered stone rise from the seafloor, their spires and domes home to schools of colorful fish and undulating anemones.\n\nStreets and plazas, once bustling with life, are now silent corridors through which curious sea turtles and shy octopi navigate. Statues of forgotten heroes stand sentinel, their features softened by the relentless caress of the currents.\n\nIn the grand palace at the city's center, treasures beyond imagination still glitter in the filtered sunlight – gold and jewels untouched by mortal hands for centuries. But they are guarded by ancient magic and fearsome creatures of the deep.\n\nThose who dare to explore the Sunken City may discover the secrets of a lost civilization, decipher the runes etched into waterlogged tablets, or awaken slumbering powers best left undisturbed. The weight of history and the pressure of the deep make every moment in this submerged wonder a test of courage and will.",
		},
		{
			"title":    "The Celestial Citadel",
			"desc":     "A floating fortress in the sky",
			"longDesc": "High above the world, suspended among the clouds, floats the Celestial Citadel. This majestic fortress of gleaming silver and gold seems to defy gravity, its impossible architecture a testament to the power of ancient sky-magic.\n\nSpires of crystal and adamantine pierce the heavens, while great arches and flying buttresses support expansive courtyards and gardens. The air is thin but invigorating, filled with the scent of rare floating flowers and the distant song of wind spirits.\n\nCelestial beings of light and air walk the hallways, their feet barely touching the ground. In the great observatory, star-charts of unimaginable detail map not just this world, but countless others across the cosmos.\n\nAt the very peak of the highest tower sits the Throne of Winds, from which it is said one can command the very forces of the sky. But the path to this seat of power is fraught with peril – treacherous updrafts, labyrinthine corridors that shift with the changing winds, and tests of worthiness that challenge both body and spirit.\n\nThose who seek the Celestial Citadel must not only find a way to reach these lofty heights but also prove themselves worthy of the wonders and dangers that await in this realm where earth meets sky.",
		},
	}
	items := make([]components.Item, len(itemData))
	for i, item := range itemData {
		items[i] = components.NewItem(item["title"], item["desc"], item["longDesc"])
	}

	component := components.NewListViewportModel(items)

	p := tea.NewProgram(model{component: component}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
