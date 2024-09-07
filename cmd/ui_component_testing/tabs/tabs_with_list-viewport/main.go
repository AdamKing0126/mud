package main

import (
	"fmt"
	"log/slog"
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
	title := titleStyle.Render("Tabs With List and Viewport")
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
	logFile, err := os.OpenFile("tabs_testing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{}))

	tabData := []struct {
		name  string
		items []components.Item
	}{
		{
			name: "Locations",
			items: []components.Item{
				components.NewItem("1", "area", "The Enchanted Forest", "A mystical woodland realm", `The Enchanted Forest is a realm of wonder and magic, where ancient trees whisper secrets to those who listen. Sunlight filters through a canopy of emerald leaves, casting dappled shadows on a carpet of moss and wildflowers. Ethereal mists curl around gnarled roots, concealing hidden pathways and forgotten ruins.

Mysterious creatures flit between the trees – iridescent fairies, mischievous sprites, and wise old talking animals. The air is thick with the scent of pine and wild herbs, and the sound of tinkling laughter echoes from unseen sources.

At the heart of the forest stands the ancient Tree of Life, its branches reaching towards the heavens and its roots delving deep into the earth. Legend says that those who touch its bark can glimpse visions of past and future.

Brave adventurers who enter the Enchanted Forest may find themselves facing magical trials, solving riddles posed by enigmatic beings, or stumbling upon portals to other realms. But beware – the forest has a mind of its own, and paths have a habit of shifting when you're not looking.`),
				components.NewItem("2", "area", "The Sunken City", "An underwater metropolis", `Beneath the rolling waves of the ocean lies the Sunken City, a marvel of ancient engineering and magic. Once a thriving surface metropolis, it now rests on the ocean floor, protected from the crushing depths by a shimmering dome of arcane energy.

Streets paved with luminescent pearls wind between towering coral-encrusted spires and grand plazas where schools of colorful fish dart like living confetti. The city's inhabitants – a mix of merfolk, adapted humans, and exotic sea creatures – go about their daily lives in this aquatic wonderland.

At the city's center stands the majestic Palace of a Thousand Currents, home to the Coral Throne and the mysterious Sea Queen who rules over this submerged realm. The palace's halls are said to hold treasures beyond imagination, guarded by fearsome leviathans and cunning water elementals.

Explorers brave enough to visit the Sunken City must contend not only with the dangers of the deep but also with the city's complex politics and ancient magics. But for those who succeed, the rewards – whether knowledge, treasure, or alliance – can be as vast as the ocean itself.`),
				components.NewItem("3", "area", "The Celestial Citadel", "A floating fortress in the sky", `High above the world, suspended among the clouds, floats the Celestial Citadel – a breathtaking fortress of gleaming silver spires and gossamer bridges. This bastion of the heavens is home to a society of skyfarers, astronomers, and elemental air mages who have turned their backs on the earth below.

The citadel is a masterpiece of magical architecture, with floating gardens, crystalline observatories, and grand halls where the floor is nothing but swirling mist. Gravity is a suggestion here, with residents and visitors alike able to drift between levels with a thought.

At the pinnacle of the highest tower sits the Astral Nexus, a powerful artifact that allows the citadel's rulers to chart the movements of celestial bodies and peer into other planes of existence. It is said that on clear nights, one can see the strands of fate themselves from the Nexus chamber.

While the Celestial Citadel may seem like a paradise, it faces constant threats from sky pirates, jealous dragons, and the very storms it drifts through. Visitors must prove their worth and intentions before being granted entry to this wondrous realm among the stars.`),
			},
		},
		{
			name: "Characters",
			items: []components.Item{
				components.NewItem("1", "character", "Elara the Elven Mage", "Master of arcane arts", `Elara is a centuries-old elven mage, renowned for her mastery of elemental magic and her insatiable thirst for knowledge. With hair like spun silver and eyes that shimmer with arcane power, she cuts an imposing figure in her flowing robes embroidered with mystical sigils.

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
			},
		},
		{
			name: "Artifacts",
			items: []components.Item{
				components.NewItem("1", "item", "The Orb of Destiny", "A crystal ball of immense power", `The Orb of Destiny is a legendary artifact that grants visions of possible futures to those strong-willed enough to use it. About the size of a human skull, the Orb appears to be made of flawless crystal, within which swirl mists of every color imaginable.

Created eons ago by a conclave of the most powerful seers and chronomancers in history, the Orb was intended to be a tool for guiding civilization towards a golden age. However, its creators soon discovered that the future is not easily shaped, and that knowledge of what may come can be as much a curse as a blessing.

To use the Orb, one must enter a deep meditative state while gazing into its depths. The mists within will part, revealing vivid visions of potential futures. These visions are not set in stone, but rather show the most likely outcomes based on current circumstances and the choices of individuals.

The power of the Orb comes with great risk. Unprepared or weak-minded users can become lost in the visions, their consciousness trapped within the Orb for eternity. Even for those strong enough to withstand its power, the knowledge gained can be maddening, showing both wondrous possibilities and terrible catastrophes that may come to pass.

Throughout history, the Orb has passed through the hands of many rulers, mages, and adventurers. Some have used its power wisely, averting disasters and guiding their people towards prosperity. Others have been driven to despair or madness by the weight of the knowledge it provides.

The current whereabouts of the Orb are unknown, but rumors persist of its appearance in times of great crisis. Those who seek it must ask themselves: are they prepared to bear the burden of knowing what the future may hold?`),
				components.NewItem("1", "item", "The Sword of the Dawn", "A blade of pure light", `Forged from a fallen star and imbued with the essence of the rising sun, the Sword of the Dawn is a weapon of unparalleled radiance and power. The blade itself appears to be made of pure, solidified light, its edge so keen it is said to be able to cut through shadows themselves.

The hilt of the sword is crafted from orichalcum, a legendary metal of celestial origin, and is adorned with intricate engravings depicting the eternal dance of light and darkness. When wielded, the sword emits a soft, warm glow that intensifies to a brilliant radiance in the presence of evil or dark magic.

Legend has it that the Sword of the Dawn was created by the gods themselves as a gift to mortalkind, meant to be a beacon of hope in times of darkness. Its first wielder was a humble farmer who rose to become a great hero, using the sword's power to banish an encroaching darkness that threatened to engulf the world.

The sword's most potent ability is its power to banish darkness in all its forms. It can dispel illusions, break curses, and even harm creatures of shadow and night that are normally impervious to physical attacks. In the hands of a pure-hearted wielder, its light can heal wounds and inspire courage in allies.

However, the Sword of the Dawn is not without its dangers. Its power is tied to the cycle of day and night, waxing and waning with the sun. At high noon, its abilities are at their peak, but as night falls, the sword's glow dims, and it becomes little more than an ordinary blade. Moreover, those of evil heart who attempt to wield the sword find its touch painfully hot, and prolonged contact can even ignite them in purifying flames.

Throughout history, the Sword of the Dawn has appeared in times of great darkness, always finding its way into the hands of those destined to push back the shadows. Its current whereabouts are unknown, but many believe that when the world faces its darkest hour, the Sword of the Dawn will once again emerge to light the way towards salvation.`),
				components.NewItem("3", "item", "The Cloak of Shadows", "A garment of pure darkness", `Woven from the essence of shadow itself, the Cloak of Shadows is an artifact of immense power that grants its wearer mastery over darkness and stealth. At first glance, it appears to be a simple hooded cloak of the deepest black. However, closer inspection reveals that it seems to absorb light, with faint wisps of shadow constantly swirling across its surface.

The origins of the Cloak are shrouded in mystery and conflicting legends. Some say it was created by the first thief to ever successfully steal from the gods. Others claim it was woven by a forgotten goddess of night and secrets as a gift to her most devoted followers. Whatever the truth, the Cloak has been a coveted item among thieves, spies, and assassins for centuries.

When worn, the Cloak grants its user a range of shadow-based abilities. The most basic of these is near-perfect invisibility in any area of shade or darkness. Beyond this, skilled users can shape shadows into solid forms, create areas of magical darkness, and even step through shadows to teleport short distances.

The true extent of the Cloak's powers is said to be even greater. Masters of the Cloak are rumored to be able to become one with the shadows, transforming their body into living darkness. Some tales even speak of wearers using the Cloak to travel to the Plane of Shadow, a dark mirror of the material world.

However, the Cloak of Shadows is not without its risks. Prolonged use can have a corrupting influence, causing the wearer to become increasingly secretive, paranoid, and drawn to darkness. There are stories of wearers who became so enamored with the shadows that they forgot their own identities, eventually fading away into the darkness forever.

The Cloak of Shadows has changed hands countless times throughout history, often through theft or assassination. Its current owner, if it has one, is unknown – which is exactly how the wearer would prefer it. Those who seek the Cloak must be prepared not only to find it but to resist the temptation of the absolute secrecy and power it offers.`),
			},
		},
	}

	var tabContent []components.Component
	var tabNames []string

	for _, tab := range tabData {
		tabNames = append(tabNames, tab.name)
		tabContent = append(tabContent, components.NewListViewportModel(tab.items, logger))
	}

	tabsComponent := components.NewTabsModel(tabNames, tabContent, logger)

	p := tea.NewProgram(model{component: tabsComponent}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
