package mobs_test

import (
	"mud/interfaces"
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

type MockAction struct{}

func (a *MockAction) GetName() string {
	return "Mock Action"
}

func (a *MockAction) GetDescription() string {
	return "This is a mock action"
}

func (a *MockAction) GetAttackBonus() int32 {
	return 5
}

func (a *MockAction) GetDamageDice() string {
	return "1d6"
}

func (a *MockAction) GetDamageBonus() int32 {
	return 2
}

func TestExecuteAction(t *testing.T) {
	mob := &mobs.Mob{
		Name:    "Monster",
		Actions: []interfaces.MobAction{&MockAction{}},
		// Set other fields as necessary...
	}

	opponent := &MockOpponent{}

	mob.ExecuteAction(opponent)

	// Add assertions here to verify the behavior of ExecuteAction.
	// For example, you might want to check that the opponent's HP decreased,
	// or that the correct action was chosen, etc.
}
