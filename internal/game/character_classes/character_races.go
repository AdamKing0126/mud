package character_classes

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/adamking0126/mud/pkg/database"
)

type CharacterRaces []CharacterRace

type Asi struct {
	Attributes []string `json:"attributes"`
	Value      int      `json:"value"`
}

// for now, just implement the bare minimum for race
type CharacterRace struct {
	Name               string `db:"name"`
	Slug               string `db:"slug"`
	SubRaceName        string `db:"subrace_name"`
	SubRaceSlug        string `db:"subrace_slug"`
	SubRaceDescription string `db:"subrace_description"`
	Description        string `db:"description"`
	ASIData            string `db:"asi"` //ability score improvements
	// TODO: Need to pull this one in, because some races have "Any" and "Other" which are described in this field
	// Eh, maybe just say, "fuck it"
	// ASIDescription     string `db:"asi_description`
	ASI []Asi
}

func (c CharacterRaces) SubRacesFor(slug string) CharacterRaces {
	var subraces CharacterRaces
	for _, characterRace := range c {
		if characterRace.Slug == slug {
			subraces = append(subraces, characterRace)
		}
	}
	return subraces
}

func (c CharacterRaces) SubRaceNamesFor(slug string) []string {
	var subraceNames []string
	for _, characterRace := range c {
		if characterRace.Slug == slug {
			subraceNames = append(subraceNames, characterRace.SubRaceName)
		}
	}
	sort.Strings(subraceNames)
	return subraceNames
}

func (c CharacterRaces) GetCharacterRaceBySubRaceName(subRaceName string) *CharacterRace {
	for _, characterRace := range c {
		if characterRace.SubRaceName == subRaceName {
			return &characterRace
		}
	}
	return nil
}

func (c CharacterRaces) GetCharacterRaceByName(raceName string) *CharacterRace {
	for _, characterRace := range c {
		if characterRace.Name == raceName {
			return &characterRace
		}
	}
	return nil
}

func GetCharacterRaceList(ctx context.Context, db database.DB, raceSlug string, subRaceSlug string) (CharacterRaces, error) {
	const baseQuery = `SELECT name, slug, subrace_name, subrace_slug, description, asi, subrace_description FROM character_races`
	var query string
	args := []interface{}{}

	if raceSlug != "" {
		if subRaceSlug != "" {
			query = baseQuery + " WHERE slug = ? AND subrace_slug = ?;"
			args = append(args, raceSlug, subRaceSlug)
		} else {
			query = baseQuery + " WHERE slug = ?;"
			args = append(args, raceSlug)
		}
	} else {
		query = baseQuery + " ORDER BY slug, subrace_slug;"
	}

	characterRaces := []CharacterRace{}
	err := db.Select(ctx, &characterRaces, query, args...)

	for i := range characterRaces {
		err := json.Unmarshal([]byte(characterRaces[i].ASIData), &characterRaces[i].ASI)
		if err != nil {
			return nil, err
		}
	}

	return characterRaces, err
}
