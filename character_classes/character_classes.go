package character_classes

import (
	"sort"

	"github.com/jmoiron/sqlx"
)

type CharacterClass struct {
	Name                    string `db:"name"`
	HitDice                 string `db:"hit_dice"`
	HPAtFirstLevel          int    `db:"hp_at_first_level"`
	HPModifier              string `db:"hp_modifier"`
	SavingThrowCharisma     bool   `db:"saving_throw_charisma"`
	SavingThrowConstitution bool   `db:"saving_throw_constitution"`
	SavingThrowDexterity    bool   `db:"saving_throw_dexterity"`
	SavingThrowIntelligence bool   `db:"saving_throw_intelligence"`
	SavingThrowStrength     bool   `db:"saving_throw_strength"`
	SavingThrowWisdom       bool   `db:"saving_throw_wisdom"`
	Slug                    string `db:"slug"`
	ArchetypeSlug           string `db:"archetype_slug"`
	ArchetypeName           string `db:"archetype_name"`
	ArchetypeDescription    string `db:"archetype_description"`
}

type CharacterClasses []CharacterClass

func (c CharacterClasses) ArchetypesFor(slug string) []CharacterClass {
	var archetypes []CharacterClass
	for _, characterClass := range c {
		if characterClass.Slug == slug {
			archetypes = append(archetypes, characterClass)
		}
	}
	return archetypes
}

func (c CharacterClasses) ArchetypeNamesFor(slug string) []string {
	var archetypeNames []string
	for _, characterClass := range c {
		if characterClass.Slug == slug {
			archetypeNames = append(archetypeNames, characterClass.ArchetypeName)
		}
	}
	sort.Strings(archetypeNames)
	return archetypeNames
}

func (c CharacterClasses) GetCharacterClassByArchetypeName(archetypeName string) *CharacterClass {
	for _, characterClass := range c {
		if characterClass.ArchetypeName == archetypeName {
			return &characterClass
		}
	}
	return nil
}

func (c CharacterClasses) GetCharacterClassByName(className string) *CharacterClass {
	for _, characterClass := range c {
		if characterClass.Name == className {
			return &characterClass
		}
	}
	return nil
}

func (c CharacterClasses) GetCharacterClassByArchetypeSlug(archetypeSlug string) *CharacterClass {
	for _, characterClass := range c {
		if characterClass.ArchetypeSlug == archetypeSlug {
			return &characterClass
		}
	}
	return nil
}

func GetCharacterClassList(db *sqlx.DB, archetype_slug string) (CharacterClasses, error) {
	const baseQuery = `SELECT name, hit_dice, hp_at_first_level, hp_modifier, saving_throw_charisma, saving_throw_constitution, saving_throw_dexterity, saving_throw_intelligence, saving_throw_strength, saving_throw_wisdom, slug, archetype_slug, archetype_name, archetype_description FROM character_classes`
	var query string
	args := []interface{}{}

	if archetype_slug != "" {
		query = baseQuery + " WHERE archetype_slug = ?;"
		args = append(args, archetype_slug)
	} else {
		query = baseQuery + " ORDER BY slug, archetype_slug;"
	}

	characterClasses := []CharacterClass{}
	err := db.Select(&characterClasses, query, args...)
	return characterClasses, err
}
