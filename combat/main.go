package combat

import (
	"sort"
)

type Combat struct {
	Aggressors []Combatant
	Defenders  []Combatant
	TurnOrder  []Combatant
}

func NewCombat(aggressors []Combatant, defenders []Combatant) Combat {
	return Combat{Aggressors: aggressors, Defenders: defenders}
}

func (c *Combat) AddAggressor(aggressor Combatant) {
	c.Aggressors = append(c.Aggressors, aggressor)
	// if the combatant joins the combat after it has begun, they
	// go to the end of the TurnOrder
	if len(c.TurnOrder) > 0 {
		c.TurnOrder = append(c.TurnOrder, aggressor)
	}
}

func (c *Combat) AddDefender(defender Combatant) {
	// if the combatant joins the combat after it has begun, they
	// go to the end of the TurnOrder
	c.Defenders = append(c.Defenders, defender)
	if len(c.TurnOrder) > 0 {
		c.TurnOrder = append(c.TurnOrder, defender)
	}
}

func (c *Combat) RollInitiative() {
	type InitiativeCombatant struct {
		Combatant      Combatant
		InitiativeRoll int32
	}

	combatants := []InitiativeCombatant{}

	for idx := range c.Aggressors {
		initiativeRoll := c.Aggressors[idx].RollInitiative()
		combatants = append(combatants, InitiativeCombatant{
			Combatant: c.Aggressors[idx], InitiativeRoll: initiativeRoll})
	}

	for idx := range c.Defenders {
		initiativeRoll := c.Defenders[idx].RollInitiative()
		combatants = append(combatants, InitiativeCombatant{
			Combatant:      c.Defenders[idx],
			InitiativeRoll: initiativeRoll})
	}

	sort.Slice(combatants, func(i, j int) bool {
		return combatants[i].InitiativeRoll > combatants[j].InitiativeRoll
	})

	for idx := range combatants {
		c.TurnOrder = append(c.TurnOrder, combatants[idx].Combatant)
	}
}
