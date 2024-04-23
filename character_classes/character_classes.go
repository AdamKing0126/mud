package character_classes

import (
	"github.com/jmoiron/sqlx"
)

type CharacterClass struct {
	Name                 string `db:"name"`
	Slug                 string `db:"slug"`
	ArchetypeSlug        string `db:"archetype_slug"`
	ArchetypeName        string `db:"archetype_name"`
	ArchetypeDescription string `db:"archetype_description"`
}

func GetCharacterClassList(db *sqlx.DB) ([]CharacterClass, error) {
	characterClasses := []CharacterClass{}

	query := `SELECT name, slug, archetype_slug, archetype_name, archetype_description FROM character_classes ORDER BY slug, archetype_slug;`
	err := db.Select(&characterClasses, query)

	return characterClasses, err
}
