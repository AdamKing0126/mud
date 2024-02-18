package interfaces

type Combat interface {
	GetAggressors() []Combatant
	GetDefenders() []Combatant
	AddAggressor(Combatant)
	AddDefender(Combatant)
	GetTurnOrder() []Combatant
}
