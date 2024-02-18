package mobs

import (
	"fmt"
	"log"
	"mud/dice"
	"mud/interfaces"
	"regexp"
	"strconv"
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
	RNG                   RNG     `db:"-"`
	Actions               []interfaces.MobAction
}

type RNG interface {
	Intn(n int) int
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
	if strings.ToLower(attack.GetName()) == "multiattack" {
		if diceRoll >= opponent.GetArmorClass() {
			fmt.Printf("MultiAttack! \"%s\"\nDice Roll: %d, Target AC: %d - SUCCESS\n", attack.GetDescription(), diceRoll, opponent.GetArmorClass())
			return true
		} else {
			fmt.Printf("MultiAttack! Dice Roll: %d, Target AC: %d - FAILED\n", diceRoll, opponent.GetArmorClass())
			return false
		}
	}
	attackBonus := attack.GetAttackBonus()
	hit := diceRoll+attackBonus >= opponent.GetArmorClass()
	fmt.Printf("Dice Roll: %d, Mob Attack Bonus: %d, Target AC: %d", diceRoll, attackBonus, opponent.GetArmorClass())
	if diceRoll == 20 {
		fmt.Println(" - Automatic HIT!")
		return true
	} else if diceRoll == 1 {
		fmt.Println(" - Automatic MISS!")
		return false
	}
	if hit {
		fmt.Printf(" - HIT\n")
	} else {
		fmt.Printf(" - MISS\n")
	}
	return hit
}

type Opponent struct {
	ArmorClass int32
}

func (mob *Mob) ExecuteRegularAttack(opponent interfaces.Opponent, attack interfaces.MobAction) {
	if mob.AttackRoll(opponent, attack) {
		// Regular expression to match damage types:
		// re := regexp.MustCompile(`(\w+) damage`)
		// re := regexp.MustCompile(`\(\d+d\d+\+\d+\) (\w+)`)
		re := regexp.MustCompile(`\(\d+d\d+(\s*\+\s*\d+)?\) (\w+)`)

		// Find all damage types in the desc field:
		matches := re.FindAllStringSubmatch(attack.GetDescription(), -1)

		// Split the damage_dice field to get the dice for each damage type:
		damageDice := strings.Split(attack.GetDamageDice(), "+")

		// Create a slice to hold the resulting damage instructions:
		result := make([][]string, len(matches))

		// Iterate over the matches and dice and add them to the result slice:
		if len(matches) == 1 {
			result[0] = []string{matches[0][len(matches[0])-1], strconv.Itoa(int(dice.DiceRoll(attack.GetDamageDice())))}
		} else {
			for i, match := range matches {
				result[i] = []string{match[len(match)-1], strconv.Itoa(int(dice.DiceRoll(damageDice[i])))}
			}
		}
		// attackDamage := dice.DiceRoll(attack.GetDamageDice())
		// fmt.Printf("%s attacks! %s does %d damage\n", mob.GetName(), attack.GetName(), attackDamage)
		// need to decrease the opponent's health.
		if len(result) == 1 {
			fmt.Printf("%s's %s does %s %s damage\n", mob.GetName(), attack.GetName(), result[0][1], result[0][0])
		} else if len(result) == 2 {
			var damageResults []string
			for idx := range result {
				damageResults = append(damageResults, fmt.Sprintf("%s %s damage", result[idx][1], result[idx][0]))
			}
			fmt.Printf("%s's %s does %s\n", mob.GetName(), attack.GetName(), strings.Join(damageResults, " and "))
		} else if len(result) > 2 {
			var damageResults []string
			for idx := range result {
				damageResults = append(damageResults, fmt.Sprintf("%s %s damage", result[idx][1], result[idx][0]))
			}
			last := damageResults[len(damageResults)-1]
			damageResults = damageResults[:len(damageResults)-1]

			fmt.Printf("%s's %s does %s\n", mob.GetName(), attack.GetName(), strings.Join(damageResults, ", ")+", and "+last)

		} else {
			fmt.Println("whoops!")
		}
		// TODO Adam need to decrement Opponent hitpoints
		fmt.Printf("yo!")
	}
}

func (mob *Mob) ExecuteAction(opponent interfaces.Opponent) {
	// Mob has a bunch of actions, need to pick one
	// "regular" actions are ones which have DamageDice - put those into a bucket
	actions := mob.GetActions()
	regularAttacks := getRegularAttacks(actions)
	multiAttack := getMultiAttack(actions)

	if multiAttack != nil && mob.AttackRoll(opponent, multiAttack) {
		attacks := getAttacksForMultiAttack(multiAttack.GetDescription(), regularAttacks)
		for idx := range attacks {
			mob.ExecuteRegularAttack(opponent, attacks[idx])
		}
	} else {
		// if the mob fails a roll to make a multiattack OR if the mob
		// does not have multiattack capability, select a regular attack at
		// random and execute it
		index := mob.RNG.Intn(len(regularAttacks))
		attack := regularAttacks[index]
		mob.ExecuteRegularAttack(opponent, attack)
	}
}

func getMultiAttack(mobActions []interfaces.MobAction) interfaces.MobAction {
	for idx := range mobActions {
		if strings.ToLower(mobActions[idx].GetName()) == "multiattack" {
			description := strings.ToLower(mobActions[idx].GetDescription())
			description = strings.ReplaceAll(description, "one", "1")
			description = strings.ReplaceAll(description, "two", "2")
			description = strings.ReplaceAll(description, "three", "3")
			description = strings.ReplaceAll(description, "four", "4")
			description = strings.ReplaceAll(description, "five", "5")
			description = strings.ReplaceAll(description, "six", "6")
			description = strings.ReplaceAll(description, "seven", "7")
			description = strings.ReplaceAll(description, "eight", "8")
			description = strings.ReplaceAll(description, "nine", "9")
			description = strings.ReplaceAll(description, "ten", "10")
			mobActions[idx].SetDescription(description)
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

func getAttacksForMultiAttack(multiAttackDescription string, regularAttacks []interfaces.MobAction) []interfaces.MobAction {
	re := regexp.MustCompile(`(\d+)\s+([a-z]\w*)`)
	matches := re.FindAllStringSubmatch(multiAttackDescription, -1)

	actionCounts := make(map[string]int32)
	for _, match := range matches {
		count, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatalf("error converting string to int: %v", err)
		}
		action := match[2]
		actionCounts[action] = int32(count)
	}

	result := []interfaces.MobAction{}
	for _, attack := range regularAttacks {
		count, ok := actionCounts[strings.ToLower(attack.GetName())]
		if ok {
			for i := int32(0); i < count; i++ {
				result = append(result, attack)
			}
		}
	}
	return result
}
