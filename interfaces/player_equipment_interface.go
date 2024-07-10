package interfaces

type PlayerEquipment interface {
	GetUUID() string
	GetPlayerUUID() string
	GetHead() *EquippedItem
	GetNeck() *EquippedItem
	GetChest() *EquippedItem
	GetArms() *EquippedItem
	GetHands() *EquippedItem
	GetDominantHand() *EquippedItem
	GetOffHand() *EquippedItem
	GetLegs() *EquippedItem
	GetFeet() *EquippedItem
}
