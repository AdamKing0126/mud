package players

import (
	"fmt"
	"mud/display"
	"mud/interfaces"
	"net"
	"reflect"

	"github.com/jmoiron/sqlx"
)

func NewPlayer(conn net.Conn) *Player {
	return &Player{Conn: conn}
}

type Player struct {
	UUID            string
	Name            string
	RoomUUID        string
	Room            interfaces.Room
	AreaUUID        string
	Area            interfaces.Area
	Health          int
	HealthMax       int
	Mana            int
	ManaMax         int
	Movement        int
	MovementMax     int
	Conn            net.Conn
	Commands        []string
	ColorProfile    ColorProfile
	LoggedIn        bool
	Password        string
	PlayerAbilities PlayerAbilities
	Equipment       PlayerEquipment
	Inventory       []interfaces.Item
}

func (player *Player) AddItem(db *sqlx.DB, item interfaces.Item) error {
	err := item.SetLocation(db, player.UUID, "")
	if err != nil {
		return err
	}
	player.SetInventory(append(player.GetInventory(), item))
	return nil
}

func (player *Player) RemoveItem(item interfaces.Item) error {
	itemIndex := -1
	for idx := range player.GetInventory() {
		if player.Inventory[idx] == item {
			itemIndex = idx
			break
		}
	}
	if itemIndex == -1 {
		return fmt.Errorf("item %s is not found in player %s inventory", item.GetUUID(), player.GetUUID())
	}
	player.Inventory = append(player.Inventory[:itemIndex], player.Inventory[itemIndex+1:]...)
	return nil
}

func (player *Player) Regen(db *sqlx.DB) error {
	healthRegen := calculateHealthRegen(player)
	manaRegen := calculateManaRegen(player)
	movementRegen := calculateMovementRegen(player)

	player.Health = int(float64(player.Health) * healthRegen)
	if player.Health > player.HealthMax {
		player.Health = player.HealthMax
	}

	player.Mana = int(float64(player.Mana) * manaRegen)
	if player.Mana > player.ManaMax {
		player.Mana = player.ManaMax
	}

	player.Movement = int(float64(player.Movement) * movementRegen)
	if player.Movement > player.MovementMax {
		player.Movement = player.MovementMax
	}

	_, err := db.Exec("UPDATE players SET health = ?, mana = ?, movement = ? WHERE uuid = ?", player.Health, player.Mana, player.Movement, player.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (player *Player) Remove(db *sqlx.DB, itemName string) {
	if player.Equipment.DominantHand.GetName() == itemName {
		equippedItem := player.Equipment.DominantHand
		player.AddItem(db, equippedItem)

		player.Equipment.DominantHand = nil
		queryString := fmt.Sprintf("UPDATE player_equipments SET DominantHand = '' WHERE player_uuid = '%s'", player.GetUUID())
		_, err := db.Exec(queryString)
		if err != nil {
			fmt.Printf("error setting player_equipment to nil: %v", err)
		}

		display.PrintWithColor(player, fmt.Sprintf("You remove %s.\n", equippedItem.GetName()), "reset")
	}
}

func (player *Player) Equip(db *sqlx.DB, item interfaces.Item) bool {
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
	rows, err := db.Query(queryString, player.GetUUID())
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
			_, err = db.Exec(queryString, itemUUID, player.GetUUID())
			if err != nil {
				fmt.Printf("error inserting into player_equipments: %v", err)
				return false
			}
			if columns[idx] == "DominantHand" {
				equippedItem := &PlayerEquippedItem{
					Item:         item,
					EquippedSlot: columns[idx],
				}
				player.Equipment.DominantHand = equippedItem
			}
			return true
		}
	}
	return true
}

func printEquipmentElement(player interfaces.Player, partName string, getterFunc func() interfaces.EquippedItem) {
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

	equipment := player.GetEquipment()
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
