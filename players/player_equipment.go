package players

import "mud/interfaces"

type PlayerEquipment struct {
	UUID         string
	PlayerUUID   string
	Head         interfaces.ItemInterface
	Neck         interfaces.ItemInterface
	Chest        interfaces.ItemInterface
	Arms         interfaces.ItemInterface
	Hands        interfaces.ItemInterface
	DominantHand interfaces.ItemInterface
	OffHand      interfaces.ItemInterface
	Legs         interfaces.ItemInterface
	Feet         interfaces.ItemInterface
}
