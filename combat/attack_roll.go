package combat

import (
	"github.com/adamking0126/mud/utilities"
)

func AttackRoll(combatant Combatant, opponent Combatant) bool {
	abilities := combatant.GetAbilities()
	attackModifier := abilities.GetAttackModifier("ranged")

	d20Roll := utilities.DiceRoll("1d20")
	if d20Roll == 1 {
		return false
	} else if d20Roll == 20 {
		return true
	}

	return attackModifier+d20Roll >= opponent.GetArmorClass()
}
