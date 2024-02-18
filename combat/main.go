package combat

import (
	"mud/interfaces"
	"sort"
)

type Combat struct {
	Aggressors []interfaces.Combatant
	Defenders  []interfaces.Combatant
	TurnOrder  []interfaces.Combatant
}

func NewCombat(aggressors []interfaces.Combatant, defenders []interfaces.Combatant) Combat {
	return Combat{Aggressors: aggressors, Defenders: defenders}
}

func (c *Combat) GetAggressors() []interfaces.Combatant {
	return c.Aggressors
}

func (c *Combat) GetDefenders() []interfaces.Combatant {
	return c.Defenders
}

func (c *Combat) AddAggressor(aggressor interfaces.Combatant) {
	c.Aggressors = append(c.Aggressors, aggressor)
	// if the combatant joins the combat after it has begun, they
	// go to the end of the TurnOrder
	if len(c.TurnOrder) > 0 {
		c.TurnOrder = append(c.TurnOrder, aggressor)
	}
}

func (c *Combat) AddDefender(defender interfaces.Combatant) {
	// if the combatant joins the combat after it has begun, they
	// go to the end of the TurnOrder
	c.Defenders = append(c.Defenders, defender)
	if len(c.TurnOrder) > 0 {
		c.TurnOrder = append(c.TurnOrder, defender)
	}
}

func (c *Combat) GetTurnOrder() []interfaces.Combatant {
	return c.TurnOrder
}

func (c *Combat) RollInitiative() {
	type InitiativeCombatant struct {
		Combatant      interfaces.Combatant
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
