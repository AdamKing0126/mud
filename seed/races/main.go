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
	ASI         string `json:"asi"`
}

type RaceImport struct {
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	ASI          string `db:"asi"`
	Description  string `db:"description"`
	SubracesData string `db:"subraces"`
	Subraces     []SubraceImport
}

func SeedRaces(dbPath string, racesImportDbPath string) error {
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

	query := `SELECT name, slug, subraces, description, asi from race_imports;`
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
		var subraces []SubraceImport
		err = json.Unmarshal([]byte(ri.SubracesData), &subraces)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		ri.Subraces = subraces

	}
	return nil
}
