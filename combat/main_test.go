package combat

import (
	"mud/interfaces"
	"testing"
)

type TestCombatant struct {
	initiativeToReturn int32
	armorClass         int32
}

func (c *TestCombatant) RollInitiative() int32 {
	return c.initiativeToReturn
}

func (c *TestCombatant) GetArmorClass() int32 {
	return c.armorClass
}

func (c *TestCombatant) GetAbilities() interfaces.Abilities {
	return nil
}

func TestRollInitiative(t *testing.T) {
	testCases := []struct {
		aggressors, defenders, expected []interfaces.Combatant
	}{
		{
			aggressors: []interfaces.Combatant{
				&TestCombatant{initiativeToReturn: 10, armorClass: 10},
				&TestCombatant{initiativeToReturn: 19, armorClass: 10},
			},
			defenders: []interfaces.Combatant{
				&TestCombatant{initiativeToReturn: 5, armorClass: 10},
				&TestCombatant{initiativeToReturn: 20, armorClass: 10},
			},
			expected: []interfaces.Combatant{
				&TestCombatant{initiativeToReturn: 20, armorClass: 10},
				&TestCombatant{initiativeToReturn: 19, armorClass: 10},
				&TestCombatant{initiativeToReturn: 10, armorClass: 10},
				&TestCombatant{initiativeToReturn: 5, armorClass: 10},
			},
		},
	}

	for _, tc := range testCases {
		combat := NewCombat(tc.aggressors, tc.defenders)
		combat.RollInitiative()
		turnOrder := combat.GetTurnOrder()

		for i := range turnOrder {
			actual, ok := turnOrder[i].(*TestCombatant)
			if !ok {
				t.Errorf("turnOrder[%v] is not a *TestCombatant", i)
			}
			expected, ok := tc.expected[i].(*TestCombatant)
			if !ok {
				t.Errorf("expectedTurnOrder[%v] is not a *TestCombatant", i)
			}
			if *actual != *expected {
				t.Errorf("At index %v: expected %v, but got %v", i, *expected, *actual)
			}
		}
	}
}
