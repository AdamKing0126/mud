package players

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/character_classes"
	"github.com/adamking0126/mud/pkg/database"

	"github.com/charmbracelet/ssh"

	"golang.org/x/crypto/bcrypt"
)

func getPlayerInput(session ssh.Session) string {
	var input strings.Builder
	for {
		b := make([]byte, 1)
		_, err := session.Read(b)
		if err != nil {
			log.Printf("Error reading input: %v", err)
			return ""
		}
		if b[0] == '\r' || b[0] == '\n' {
			fmt.Fprintln(session) // Print newline
			break
		}
		input.Write(b)
		fmt.Fprint(session, string(b)) // Echo the character
	}
	return strings.TrimSpace(input.String())
}
func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

func getCharacterClassNamesFromCharacterClassObjects(c character_classes.CharacterClasses) []string {
	characterClassNamesSet := make(map[string]bool)
	for _, characterClass := range c {
		characterClassNamesSet[characterClass.Name] = true
	}
	characterClassNames := make([]string, 0, len(characterClassNamesSet))
	for k := range characterClassNamesSet {
		characterClassNames = append(characterClassNames, k)
	}
	sort.Strings(characterClassNames)
	return characterClassNames
}

func getCharacterRaceNamesFromCharacterRaceObjects(c character_classes.CharacterRaces) []string {
	characterRaceNamesSet := make(map[string]bool)
	for _, characterRace := range c {
		characterRaceNamesSet[characterRace.Name] = true
	}
	characterRaceNames := make([]string, 0, len(characterRaceNamesSet))
	for k := range characterRaceNamesSet {
		characterRaceNames = append(characterRaceNames, k)
	}

	sort.Strings(characterRaceNames)
	return characterRaceNames
}

func selectCharacterClassAndArchetype(ctx context.Context, session ssh.Session, db database.DB, player *Player) *character_classes.CharacterClass {
	characterClasses, err := character_classes.GetCharacterClassList(ctx, db, "")
	if err != nil {
		log.Fatal(err)
	}
	characterClassNames := getCharacterClassNamesFromCharacterClassObjects(characterClasses)

	for {
		menuTitle := "Choose a Character Class"
		delimiter := "number"
		lineStyle := "double"

		menuContents := display.MenuContents{
			Title:     &menuTitle,
			Delimiter: &delimiter,
			LineStyle: &lineStyle,
			MaxWidth:  65,
			Elements:  characterClassNames,
		}
		display.PrintMenu(player, menuContents)

		display.PrintWithColor(player, fmt.Sprintf("Make a selection (1-%d, anything else to quit): ", len(characterClassNames)), "primary")
		choice := getPlayerInput(session)

		characterClassChoice, err := strconv.Atoi(choice)
		if err != nil || characterClassChoice > len(characterClassNames) || characterClassChoice < 1 {
			return nil
		}
		chosenCharacterClass := characterClasses.GetCharacterClassByName(characterClassNames[characterClassChoice-1])

		for {
			archetypes := characterClasses.ArchetypesFor(chosenCharacterClass.Slug)
			archetypeNames := characterClasses.ArchetypeNamesFor(chosenCharacterClass.Slug)

			menuContents.Elements = archetypeNames
			menuTitle = fmt.Sprintf("Select a %s Subclass", chosenCharacterClass.Name)
			display.PrintMenu(player, menuContents)
			display.PrintWithColor(player, fmt.Sprintf("Make an archetype selection to learn more (1-%d, anything else to quit): ", len(archetypeNames)), "primary")

			choice = getPlayerInput(session)
			archetypeChoice, err := strconv.Atoi(choice)
			if err != nil || archetypeChoice > len(archetypes) || archetypeChoice < 1 {
				break
			}
			archetypeName := archetypeNames[archetypeChoice-1]
			chosenCharacterClass = characterClasses.GetCharacterClassByArchetypeName(archetypeName)

			menuTitle = archetypeName
			menuContents.Delimiter = nil
			menuContents.Elements = []string{chosenCharacterClass.ArchetypeDescription}
			display.PrintMenu(player, menuContents)

			display.PrintWithColor(player, fmt.Sprintf("%s\n", chosenCharacterClass.GetSavingThrowStatement()), "primary")
			display.PrintWithColor(player, fmt.Sprintf("Would you like to select a Character Class of %s-%s? (Y/N): ", chosenCharacterClass.Name, chosenCharacterClass.ArchetypeName), "primary")
			choice = getPlayerInput(session)
			choice = strings.ToLower(choice)
			fmt.Println(choice)
			if choice == "y" {
				return chosenCharacterClass
			}
			delimiter := "number"
			menuContents.Delimiter = &delimiter
		}
	}
}

func selectRace(ctx context.Context, session ssh.Session, db database.DB, player *Player) *character_classes.CharacterRace {
	characterRaces, err := character_classes.GetCharacterRaceList(ctx, db, "", "")
	if err != nil {
		log.Fatal(err)
	}
	characterRaceNames := getCharacterRaceNamesFromCharacterRaceObjects(characterRaces)

	for {
		menuTitle := "Choose a Character Race"
		delimiter := "number"
		lineStyle := "double"

		menuContents := display.MenuContents{
			Title:     &menuTitle,
			Delimiter: &delimiter,
			LineStyle: &lineStyle,
			MaxWidth:  65,
			Elements:  characterRaceNames,
		}

		display.PrintMenu(player, menuContents)

		display.PrintWithColor(player, fmt.Sprintf("Make a selection (1-%d, anything else to quit): ", len(characterRaceNames)), "primary")
		choice := getPlayerInput(session)

		characterRaceChoice, err := strconv.Atoi(choice)
		if err != nil || characterRaceChoice > len(characterRaceNames) || characterRaceChoice < 1 {
			return nil
		}

		chosenCharacterRace := characterRaces.GetCharacterRaceByName(characterRaceNames[characterRaceChoice-1])

		for {
			subRaces := characterRaces.SubRacesFor(chosenCharacterRace.Slug)

			if len(subRaces) > 1 {
				subRaceNames := characterRaces.SubRaceNamesFor(chosenCharacterRace.Slug)

				menuContents.Elements = subRaceNames
				menuTitle = fmt.Sprintf("Select a %s Subrace", chosenCharacterRace.Name)
				display.PrintMenu(player, menuContents)
				display.PrintWithColor(player, fmt.Sprintf("Make a subrace selection to learn more (1-%d, anything else to quit): ", len(subRaceNames)), "primary")

				choice = getPlayerInput(session)
				subraceChoice, err := strconv.Atoi(choice)
				if err != nil || subraceChoice > len(subRaces) || subraceChoice < 1 {
					break
				}

				subRaceName := subRaceNames[subraceChoice-1]
				chosenCharacterRace = subRaces.GetCharacterRaceBySubRaceName(subRaceName)

				menuTitle = subRaceName
				menuContents.Delimiter = nil
				menuContents.Elements = []string{chosenCharacterRace.SubRaceDescription}
				display.PrintMenu(player, menuContents)

				display.PrintWithColor(player, fmt.Sprintf("Would you like to select a Character Race of %s-%s? (Y/N): ", chosenCharacterRace.Name, chosenCharacterRace.SubRaceName), "primary")
				choice = getPlayerInput(session)
				choice = strings.ToLower(choice)
				fmt.Println(choice)
				if choice == "y" {
					return chosenCharacterRace
				}
				delimiter := "number"
				menuContents.Delimiter = &delimiter
			} else {
				display.PrintWithColor(player, fmt.Sprintf("Would you like to select a Character Race of %s? (Y/N): ", chosenCharacterRace.Name), "primary")
				choice = getPlayerInput(session)
				choice = strings.ToLower(choice)
				fmt.Println(choice)
				if choice == "y" {
					return chosenCharacterRace
				}
				delimiter := "number"
				menuContents.Delimiter = &delimiter
				break
			}
		}
	}

}

// Handle the login process for a player.  After authentication,
// cycle through related fields to populate the `player` object:
// - ColorProfile
// - Equipment
// - etc
//
// each one of these steps results in another database query, but I thought it
// best to keep the actions atomic for now, rather than trying to build one
// huge query which has joins all over the place.
func LoginPlayer(ctx context.Context, session ssh.Session, playerService *Service) (*Player, error) {

	fmt.Fprintf(session, "Welcome! Please enter your player name: ")
	playerName := getPlayerInput(session)

	player, err := playerService.GetPlayerFromDB(ctx, playerName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(session, "Player not found.  Do you want to create a new player? (y/n): ")
			answer := getPlayerInput(session)

			if strings.ToLower(answer) == "y" {
				player, err = playerService.CreatePlayer(ctx, session, playerName)
				if err != nil {
					return nil, err
				}
				return player, nil
			} else {
				return nil, errors.New("Player does not exist")
			}
		}
		return nil, err
	}

	fmt.Fprintf(session, "Please enter your password: ")
	passwd := getPlayerInput(session)
	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(passwd))
	if err != nil {
		return nil, err
	}

	player.Session = session

	playerService.SetPlayerColorProfile(ctx, player)
	playerService.SetPlayerEquipment(ctx, player)
	playerService.SetPlayerInventory(ctx, player)
	err = playerService.SetPlayerLoggedInStatus(ctx, player, true)
	if err != nil {
		return nil, err
	}

	return player, nil
}
