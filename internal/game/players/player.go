package players

import (
	"fmt"
	"reflect"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/character_classes"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/utilities"

	"github.com/charmbracelet/ssh"

	"github.com/jmoiron/sqlx"
)

func NewPlayer(session ssh.Session) *Player {
	return &Player{Session: session}
}

type Player struct {
	UUID            string
	Name            string
	RoomUUID        string
	AreaUUID        string
	HP              int32
	HPMax           int32
	Movement        int32
	MovementMax     int32
	Session         ssh.Session
	Commands        []string
	ColorProfile    ColorProfile
	LoggedIn        bool
	Password        string
	PlayerAbilities PlayerAbilities
	Equipment       PlayerEquipment
	Inventory       []*items.Item
	CharacterClass  character_classes.CharacterClass
	Race            character_classes.CharacterRace
}

func (player *Player) GetColorProfileColor(colorUse string) string {
	return player.ColorProfile.GetColor(colorUse)
}

func (player *Player) AddItem(db *sqlx.DB, item *items.Item) error {
	err := item.SetLocation(db, player.UUID, "")
	if err != nil {
		return err
	}
	player.Inventory = append(player.Inventory, item)
	return nil
}

func (player *Player) RemoveItem(item *items.Item) error {
	itemIndex := -1
	for idx := range player.Inventory {
		if player.Inventory[idx].GetUUID() == item.UUID {
			itemIndex = idx
			break
		}
	}
	if itemIndex == -1 {
		return fmt.Errorf("item %s is not found in player %s inventory", item.GetUUID(), player.UUID)
	}
	player.Inventory = append(player.Inventory[:itemIndex], player.Inventory[itemIndex+1:]...)
	return nil
}

func (player *Player) Regen(db *sqlx.DB) error {
	healthRegen := calculateHPRegen(*player)
	movementRegen := calculateMovementRegen(*player)

	player.HP = int32(float64(player.HP) * healthRegen)
	if player.HP > player.HPMax {
		player.HP = player.HPMax
	}

	player.Movement = int32(float64(player.Movement) * movementRegen)
	if player.Movement > player.MovementMax {
		player.Movement = player.MovementMax
	}

	_, err := db.Exec("UPDATE players SET hp = ?, movement = ? WHERE uuid = ?", player.HP, player.Movement, player.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (player *Player) Remove(db *sqlx.DB, itemName string) {
	if player.Equipment.DominantHand.GetName() == itemName {
		equippedItem := player.Equipment.DominantHand
		player.AddItem(db, equippedItem.Item)

		player.Equipment.DominantHand = nil
		queryString := fmt.Sprintf("UPDATE player_equipments SET DominantHand = '' WHERE player_uuid = '%s'", player.UUID)
		_, err := db.Exec(queryString)
		if err != nil {
			fmt.Printf("error setting player_equipment to nil: %v", err)
		}

		display.PrintWithColor(player, fmt.Sprintf("You remove %s.\n", equippedItem.GetName()), "reset")
	}
}

func (player *Player) Equip(db *sqlx.DB, item *items.Item) bool {
	// get the location where the thing goes
	val := reflect.ValueOf(&player.Equipment).Elem()
	itemEquipSlots := []string{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		for _, slot := range item.GetEquipmentSlots() {
			if string(slot) == field.Name {
				itemEquipSlots = append(itemEquipSlots, field.Name)
			}
		}

	}
	fmt.Println(itemEquipSlots)

	// generate the query to retrieve columns from player_equipments table
	queryString := "SELECT "
	for i, slot := range itemEquipSlots {
		if i > 0 {
			queryString += ", "
		}
		queryString += slot
	}
	queryString += " FROM player_equipments WHERE player_uuid = ? LIMIT 1"
	rows, err := db.Query(queryString, player.UUID)
	if err != nil {
		fmt.Printf("error retrieving player equipments: %v", err)
		return false
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("error getting columns: %v", err)
		return false
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))

	for i := range columns {
		pointers[i] = &values[i]
	}

	if rows.Next() {
		err = rows.Scan(pointers...)
		if err != nil {
			fmt.Printf("error scanning row: %v", err)
			return false
		}
	}

	// iterate through the values
	// get the index of the first empty value, equip the item there and then break
	for idx, val := range values {
		if val == "" {
			itemUUID := item.GetUUID()
			queryString := "UPDATE player_equipments SET "
			queryString += columns[idx]
			queryString += " = ? WHERE player_uuid = ?"
			rows.Close()
			_, err = db.Exec(queryString, itemUUID, player.UUID)
			if err != nil {
				fmt.Printf("error inserting into player_equipments: %v", err)
				return false
			}
			if columns[idx] == "DominantHand" {
				equippedItem := EquippedItem{
					Item:         item,
					EquippedSlot: columns[idx],
				}
				player.Equipment.DominantHand = &equippedItem
			}
			return true
		}
	}
	return true
}

func printEquipmentElement(player *Player, partName string, getterFunc func() *EquippedItem) {
	part := getterFunc()
	partText := "nothing"
	if part != nil {
		partText = part.GetName()
	}
	display.PrintWithColor(player, fmt.Sprintf("\n%s: %s", partName, partText), "primary")
}

func (player *Player) DisplayEquipment() {
	display.PrintWithColor(player, "\n========================", "primary")
	display.PrintWithColor(player, "\nYour current equipment:\n", "primary")

	equipment := player.Equipment
	printEquipmentElement(player, "Head", equipment.GetHead)
	printEquipmentElement(player, "Neck", equipment.GetNeck)
	printEquipmentElement(player, "Chest", equipment.GetChest)
	printEquipmentElement(player, "Arms", equipment.GetArms)
	printEquipmentElement(player, "Hands", equipment.GetHands)
	printEquipmentElement(player, "DominantHand", equipment.GetDominantHand)
	printEquipmentElement(player, "OffHand", equipment.GetOffHand)
	printEquipmentElement(player, "Legs", equipment.GetLegs)
	printEquipmentElement(player, "Feet", equipment.GetFeet)
	display.PrintWithColor(player, "\n========================\n\n", "primary")
}

func (player *Player) RollInitiative() int32 {
	// todo, there's more than this to rolling initiative but atm I can't
	// be bothered to look it up.
	return utilities.DiceRoll("1d20")
}

func (player *Player) GetColorProfilecolor(colorUse string) string {
	// TODO this ain't right.
	return player.ColorProfile.Primary
}

func (player *Player) GetSession() ssh.Session {
	return player.Session
}
