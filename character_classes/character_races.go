package character_classes

// for now, just implement the bare minimum for race
type CharacterRace struct {
	Name        string        `db:"name"`
	Slug        string        `db:"slug"`
	SubraceName string        `db:"subrace_name"`
	SubraceSlug string        `db:"subrace_slug"`
	Description string        `db:"description"`
	ASI         []interface{} `db:"asi"` //ability score improvements
}
