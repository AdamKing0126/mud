package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SubraceImport struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ASI         []Asi  `json:"asi"`
}

type Asi struct {
	Attributes []string `json:"attributes"`
	Value      int      `json:"value"`
}

type RaceImport struct {
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Size        string `db:"size_raw"`
	ASIData     string `db:"asi"`
	ASI         []Asi  `json:"asi"`
	Description string `db:"description"`
	SubraceData string `db:"subraces"`
	Subraces    []SubraceImport
}

func SeedGroups(dbPath string, racesImportDbPath string) error {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open Sqlite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Mud database opened successfully")
	}
	defer db.Close()

	racesDB, err := sqlx.Connect("sqlite3", racesImportDbPath)
	if err != nil {
		log.Fatalf("Failed to open Class Imports database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database Class Imports: %v", err)
		}
		fmt.Println("Class Imports Database opened successfully")
	}
	defer racesDB.Close()

	query := `SELECT name, slug, size_raw, subraces, description, asi from race_imports;`
	rows, err := racesDB.Queryx(query)
	if err != nil {
		log.Fatalf("Failed to query row: %v", err)
	}
	for rows.Next() {
		var ri RaceImport
		err = rows.StructScan(&ri)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		var asi []Asi
		err = json.Unmarshal([]byte(ri.ASIData), &asi)
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		ri.ASI = asi
		var subraces []SubraceImport
		err = json.Unmarshal([]byte(ri.SubraceData), &subraces)
		if err != nil {
			fmt.Println(err)
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		ri.Subraces = subraces
		// need to write each record into the game db now.

		queryString := `INSERT INTO races
		(name, slug, size, description, asi, subrace_name, subrace_slug, subrace_description)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
		if len(subraces) == 0 {
			_, err = db.Exec(queryString, ri.Name, ri.Slug, ri.ASIData, "", "", "")
			if err != nil {
				log.Fatalf("Failed to insert into races table: %v", err)
			}
		} else {
			for _, subrace := range subraces {
				asi, err := json.Marshal(subrace.ASI)
				if err != nil {
					log.Fatalf("Failed to marshal ASI into json: %v", err)
				}

				_, err = db.Exec(queryString, ri.Name, ri.Slug, asi, subrace.Name, subrace.Slug, subrace.Description)
				if err != nil {
					log.Fatalf("Failed to insert into races table: %v", err)
				}
			}
		}
	}
	return nil
}
