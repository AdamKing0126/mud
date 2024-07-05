package players

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"mud/character_classes"
	"mud/display"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"golang.org/x/crypto/bcrypt"
)

func getPlayerInput(reader io.Reader) string {
	r := bufio.NewReader(reader)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(input)
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

func selectCharacterClassAndArchetype(conn net.Conn, db *sqlx.DB, player *Player) *character_classes.CharacterClass {
	characterClasses, err := character_classes.GetCharacterClassList(db, "")
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
		choice := getPlayerInput(conn)

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

			choice = getPlayerInput(conn)
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

			display.PrintWithColor(player, fmt.Sprintf("Would you like to select a Character Class of %s-%s? (Y/N): ", chosenCharacterClass.Name, chosenCharacterClass.ArchetypeName), "primary")
			choice = getPlayerInput(conn)
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

func selectRace(conn net.Conn, db *sqlx.DB, player *Player) *character_classes.CharacterRace {
	characterRaces, err := character_classes.GetCharacterRaceList(db, "", "")
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
		choice := getPlayerInput(conn)

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

				choice = getPlayerInput(conn)
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
				choice = getPlayerInput(conn)
				choice = strings.ToLower(choice)
				fmt.Println(choice)
				if choice == "y" {
					return chosenCharacterRace
				}
				delimiter := "number"
				menuContents.Delimiter = &delimiter
			} else {
				display.PrintWithColor(player, fmt.Sprintf("Would you like to select a Character Race of %s? (Y/N): ", chosenCharacterRace.Name), "primary")
				choice = getPlayerInput(conn)
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

func createPlayer(conn net.Conn, db *sqlx.DB, playerName string) (*Player, error) {
	player := NewPlayer(conn)
	player.Name = playerName

	fmt.Fprintf(conn, "Please enter a password you'd like to use: ")
	password := getPlayerInput(conn)
	player.Password = HashPassword(password)

	// default start point
	player.AreaUUID = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
	player.RoomUUID = "189a729d-4e40-4184-a732-e2c45c66ff46"
	player.UUID = uuid.New().String()

	// "default" light mode color profile.  Should let the user choose?
	colorProfile, err := getColorProfileFromDB(db, "2c7dfd5b-d160-42e0-accb-b77d9686dbea")
	if err != nil {
		return nil, err
	}

	chosenCharacterClass := selectCharacterClassAndArchetype(conn, db, player)
	if chosenCharacterClass == nil {
		player.Logout(db)
		return nil, nil
	}

	chosenRace := selectRace(conn, db, player)
	if chosenRace == nil {
		player.Logout(db)
		return nil, nil
	}

	player.ColorProfile = *colorProfile
	player.HP = int32(chosenCharacterClass.HPAtFirstLevel)
	player.HPMax = int32(chosenCharacterClass.HPAtFirstLevel)
	player.Movement = 100
	player.MovementMax = 100
	player.UUID = uuid.New().String()
	player.Conn = conn
	player.CharacterClass = *chosenCharacterClass
	player.Race = *chosenRace

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec("INSERT INTO players (uuid, character_class, race, subrace, name, area, room, hp, hp_max, movement, movement_max, color_profile, password, logged_in) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		player.UUID, player.CharacterClass.ArchetypeSlug, player.Race.Slug, player.Race.SubRaceSlug, player.Name, player.AreaUUID, player.RoomUUID, player.HP, player.HPMax, player.Movement, player.MovementMax, player.ColorProfile.GetUUID(), player.Password, true)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert player: %v", err)
	}

	// TODO: everything below this line is probably junk.  Need to build player character based off choices made above.

	_, err = tx.Exec("INSERT INTO player_abilities (uuid, player_uuid, strength, dexterity, constitution, intelligence, wisdom, charisma) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		uuid.New(), player.UUID, 10, 10, 10, 10, 10, 10)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to set player abilities: %v", err)
	}

	_, err = tx.Exec("INSERT INTO player_equipments (uuid, player_uuid, Head, Neck, Chest, Arms, Hands, DominantHand, OffHand, Legs, Feet) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		uuid.New(), player.UUID, "", "", "", "", "", "", "", "", "", "")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to set player equipments: %v", err)
	}
	player.Equipment = *NewPlayerEquipment()

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	player.Conn = conn

	return player, nil
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
func LoginPlayer(conn net.Conn, db *sqlx.DB) (*Player, error) {

	fmt.Fprintf(conn, "Welcome! Please enter your player name: ")
	playerName := getPlayerInput(conn)

	player, err := GetPlayerFromDB(db, playerName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(conn, "Player not found.  Do you want to create a new player? (y/n): ")
			answer := getPlayerInput(conn)

			if strings.ToLower(answer) == "y" {
				player, err = createPlayer(conn, db, playerName)
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

	fmt.Fprintf(conn, "Please enter your password: ")
	passwd := getPlayerInput(conn)
	err = bcrypt.CompareHashAndPassword([]byte(player.GetHashedPassword()), []byte(passwd))
	if err != nil {
		return nil, err
	}

	player.SetConn(conn)

	err = player.GetColorProfileFromDB(db)
	if err != nil {
		return nil, err
	}

	err = player.GetEquipmentFromDB(db)
	if err != nil {
		return nil, err
	}

	err = player.GetInventoryFromDB(db)
	if err != nil {
		return nil, err
	}

	err = setPlayerLoggedInStatusInDB(db, player.GetUUID(), true)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (player *Player) Logout(db *sqlx.DB) error {
	stmt, err := db.Prepare("UPDATE players SET logged_in = FALSE WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(player.UUID)
	if err != nil {
		return err
	}

	player.Conn.Close()
	return nil
}
