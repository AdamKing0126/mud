package mobs_test

import (
	"mud/mobs"
	"testing"
)

type MockOpponent struct{}

func (o *MockOpponent) GetName() string {
	return "Mock Opponent"
}

func (o *MockOpponent) GetArmorClass() int32 {
	return 10
}

func (o *MockOpponent) GetHP() int32 {
	return 100
}

func (o *MockOpponent) DecreaseHP(amount int32) {
	// Implement this method based on your Opponent interface
}

type MockAction struct {
	Name        string
	Description string
	AttackBonus int32
	DamageDice  string
	DamageBonus int32
}

func (a *MockAction) GetName() string {
	return a.Name
}

func (a *MockAction) GetDescription() string {
	return a.Description
}

func (a *MockAction) GetAttackBonus() int32 {
	return a.AttackBonus
}

func (a *MockAction) GetDamageDice() string {
	return a.DamageDice
}

func (a *MockAction) GetDamageBonus() int32 {
	return a.DamageBonus
}

func (a *MockAction) SetDescription(desc string) {
	a.Description = desc
}

type MockRNG struct {
	IntnValue int
}

func (r *MockRNG) Intn(n int) int {
	return r.IntnValue
}

func TestExecuteAction(t *testing.T) {
	mob := &mobs.Mob{
		Name: "Monster",
		Actions: []*mobs.Action{
			&MockAction{Name: "Action1", Description: "The description for action 1 (1d6+1) bludgeoning damage", DamageDice: "1d6+1", AttackBonus: 5, DamageBonus: 5},
		},
		RNG: &MockRNG{IntnValue: 0},
	}

	opponent := &MockOpponent{}

	mob.ExecuteAction(opponent)
}

func TestExecuteActionMultiAttack(t *testing.T) {
	mob := &mobs.Mob{
		Name: "Monster",
		Actions: []*mobs.Action{
			&MockAction{Name: "Action1", Description: "The description for action 1 includes (1d6) slashing damage", DamageDice: "1d6", AttackBonus: 5, DamageBonus: 5},
			&MockAction{Name: "Action2", Description: "The description for action 2 includs (2d4) piercing", DamageDice: "1d8", AttackBonus: 3, DamageBonus: 6},
			&MockAction{Name: "Multiattack", Description: "The Monster makes one Action1 and two Action2 attacks", DamageDice: "", AttackBonus: 0, DamageBonus: 0},
		},
		RNG: &MockRNG{IntnValue: 0},
	}

	opponent := &MockOpponent{}

	mob.ExecuteAction(opponent)
}

func TestExecuteActionWithMultipleDamageTypes(t *testing.T) {
	mob := &mobs.Mob{
		Name: "Monster",
		Actions: []*mobs.Action{
			&MockAction{Name: "Action2", Description: "The description for action 2 (1d4) bludgeoning damage something (2d6) piercing.", DamageDice: "1d8+1d4", AttackBonus: 3, DamageBonus: 6},
		},
		RNG: &MockRNG{},
	}

	opponent := &MockOpponent{}

	mob.ExecuteAction(opponent)
}
