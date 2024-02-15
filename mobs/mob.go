package mobs

import (
	"fmt"
	"mud/dice"
	"mud/interfaces"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Action struct {
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"desc"`
	AttackBonus int32  `db:"attack_bonus" json:"attack_bonus,omitempty"`
	DamageDice  string `db:"damage_dice" json:"damage_dice,omitempty"`
	DamageBonus int32  `db:"damage_bonus" json:"damage_bonus,omitempty"`
}

type Mob struct {
	ID                    int64   `db:"id" mapstructure:"db"`
	AreaUUID              string  `db:"area_uuid" mapstructure:"area_uuid"`
	RoomUUID              string  `db:"room_uuid" mapstructure:"room_uuid"`
	Alignment             string  `db:"alignment" mapstructure:"alignment"`
	ArmorClass            int32   `db:"armor_class" mapstructure:"armor_class"`
	ArmorDescription      string  `db:"armor_description" mapstructure:"armor_description"`
	ChallengeRating       float64 `db:"challenge_rating" mapstructure:"challenge_rating"`
	Charisma              int32   `db:"charisma" mapstructure:"charisma"`
	CharismaSave          int32   `db:"charisma_save" mapstructure:"charisma_save"`
	ConditionImmunities   string  `db:"condition_immunities" mapstructure:"condition_immunities"`
	Constitution          int32   `db:"constitution" mapstructure:"constitution"`
	ConstitutionSave      int32   `db:"constitution_save" mapstructure:"constitution_save"`
	DamageImmunities      string  `db:"damage_immunities" mapstructure:"damage_immunities"`
	DamageResistances     string  `db:"damage_resistances" mapstructure:"damage_resistances"`
	DamageVulnerabilities string  `db:"damage_vulnerabilities" mapstructure:"damage_vulnerabilities"`
	Description           string  `db:"description" mapstructure:"description"`
	Dexterity             int32   `db:"dexterity" mapstructure:"dexterity"`
	DexteritySave         int32   `db:"dexterity_save" mapstructure:"dexterity_save"`
	Group                 string  `db:"group_name" mapstructure:"group_name"`
	HP                    int32   `db:"hp" mapstructure:"hp"`
	MaxHP                 int32   `db:"hp" mapstructure:"hp"`
	HitDice               string  `db:"hit_dice" mapstructure:"hit_dice"`
	Intelligence          int32   `db:"intelligence" mapstructure:"intelligence"`
	IntelligenceSave      int32   `db:"intelligence_save" mapstructure:"intelligence_save"`
	LegendaryDescription  string  `db:"legendary_description" mapstructure:"legendary_description"`
	Name                  string  `db:"name" mapstructure:"name"`
	Perception            int32   `db:"perception" mapstructure:"perception"`
	Senses                string  `db:"senses" mapstructure:"senses"`
	Size                  string  `db:"size" mapstructure:"size"`
	Slug                  string  `db:"slug" mapstructure:"slug"`
	Strength              int32   `db:"strength" mapstructure:"strength"`
	StrengthSave          int32   `db:"strength_save" mapstructure:"strength_save"`
	Subtype               string  `db:"subtype" mapstructure:"subtype"`
	Type                  string  `db:"type" mapstructure:"type"`
	Wisdom                int32   `db:"wisdom" mapstructure:"wisdom"`
	WisdomSave            int32   `db:"wisdom_save" mapstructure:"wisdom_save"`
	Actions               []interfaces.MobAction
}

func (mob *Mob) RollHitDice() int32 {
	// TODO: add other modifiers from feats, spells,
	// and other effects.
	return dice.DiceRoll(mob.HitDice)
}

func (mob *Mob) RollInitiative() int32 {
	// TODO: add other modifiers from feats, spells,
	// and other effects.
	return dice.DiceRoll("1d20") + int32(mob.DexteritySave)
}

func (mob *Mob) AttackRoll(opponent interfaces.Opponent, attack interfaces.MobAction) bool {
	// TODO: add/figure out attack bonus ranged vs melee etc
	diceRoll := dice.DiceRoll("1d20")
	attackBonus := attack.GetAttackBonus()
	hit := diceRoll+attackBonus >= opponent.GetArmorClass()
	fmt.Printf("Dice Roll: %d\nMob Attack Bonus: %d\nTarget AC: %d\n", diceRoll, attackBonus, opponent.GetArmorClass())
	if diceRoll == 20 {
		fmt.Println("Automatic HIT!")
		return true
	} else if diceRoll == 1 {
		fmt.Println("Automatic MISS!")
		return false
	}
	if hit {
		fmt.Printf("HIT\n")
	} else {
		fmt.Printf("MISS\n")
	}
	return hit
}

type Opponent struct {
	ArmorClass int32
}

func (mob *Mob) ExecuteAction(opponent interfaces.Opponent) {
	// Mob has a bunch of actions, need to pick one
	// "regular" actions are ones which have DamageDice - put those into a bucket
	actions := mob.GetActions()
	regularAttacks := getRegularAttacks(actions)
	multiAttack := getMultiAttack(actions)

	if multiAttack != nil {
		if mob.AttackRoll(opponent, nil) { // if the mob lands a multiattack

			// execute the multiAttacks
			// Need to split apart the attack to determine what attacks they should execute.
			// Let's keep it simple and work off of this pattern:
			// The werewolf makes two bite attacks and one claw attack
			fmt.Printf("%s\n", regularAttacks[0].GetName()) // placeholder
		}
	} else {
		// randomly(?) select one of the regular attacks
		attack := regularAttacks[0]
		if mob.AttackRoll(opponent, attack) {
			attackDamage := dice.DiceRoll(attack.GetDamageDice())
			fmt.Printf("%s attacks! %s does %d damage\n", mob.GetName(), attack.GetName(), attackDamage)
			// need to decrease the opponent's health.
		}
	}
}

func getMultiAttack(mobActions []interfaces.MobAction) interfaces.MobAction {
	for idx := range mobActions {
		if strings.ToLower(mobActions[idx].GetName()) == "multiattack" {
			return mobActions[idx]
		}
	}
	return nil
}

func getRegularAttacks(mobActions []interfaces.MobAction) []interfaces.MobAction {
	var regularAttacks []interfaces.MobAction
	for idx := range mobActions {
		if mobActions[idx].GetDamageDice() != "" {
			regularAttacks = append(regularAttacks, mobActions[idx])
		}
	}
	return regularAttacks
}
