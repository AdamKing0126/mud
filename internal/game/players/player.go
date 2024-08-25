package players

import (
	"context"
	"fmt"

	"github.com/adamking0126/mud/internal/display"
	"github.com/adamking0126/mud/internal/game/character_classes"
	"github.com/adamking0126/mud/internal/game/items"
	"github.com/adamking0126/mud/internal/utilities"
	"github.com/adamking0126/mud/pkg/database"

	"github.com/charmbracelet/ssh"
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

func (player *Player) AddItem(ctx context.Context, db database.DB, item *items.Item) error {
	err := item.SetLocation(ctx, db, player.UUID, "")
	if err != nil {
		return err
	}
	player.Inventory = append(player.Inventory, item)
	return nil
}

func (player *Player) Regen(ctx context.Context, db database.DB) error {
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

	err := db.Exec(ctx, "UPDATE players SET hp = ?, movement = ? WHERE uuid = ?", player.HP, player.Movement, player.UUID)
	if err != nil {
		return err
	}
	return nil
}

func (player *Player) Remove(ctx context.Context, db database.DB, itemName string) {
	if player.Equipment.DominantHand.GetName() == itemName {
		equippedItem := player.Equipment.DominantHand
		player.AddItem(ctx, db, equippedItem.Item)

		player.Equipment.DominantHand = nil
		queryString := fmt.Sprintf("UPDATE player_equipments SET DominantHand = '' WHERE player_uuid = '%s'", player.UUID)
		err := db.Exec(ctx, queryString)
		if err != nil {
			fmt.Printf("error setting player_equipment to nil: %v", err)
		}

		display.PrintWithColor(player, fmt.Sprintf("You remove %s.\n", equippedItem.GetName()), "reset")
	}
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
