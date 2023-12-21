package combat

import (
	"mud/dice"
	"mud/interfaces"
)

func AttackRoll(combatant interfaces.CombatantInterface, opponent interfaces.CombatantInterface) bool {
	abilities := combatant.GetAbilities()
	attackModifier := abilities.GetAttackModifier("ranged")

	d20Roll := dice.DiceRoll(1, 20)
	if d20Roll == 1 {
		return false
	} else if d20Roll == 20 {
		return true
	}

	return attackModifier+d20Roll >= opponent.GetArmorClass()
}
