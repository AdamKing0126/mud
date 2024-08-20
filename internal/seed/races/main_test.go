package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestSeedRaces(t *testing.T) {
	// housekeeping before running tests
	// need to change directories to be able to find the db files
	err := os.Chdir("../..")
	if err != nil {
		t.Fatal(err)
	}

	// delete the old test db
	dbPath := "./sql_database/test_mud.db"
	_ = os.Remove(dbPath)

	// create new db
	cmd := exec.Command("make", "create_tables", "test=true")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("some kinda error? %v", err)
	}

	raceImportDbPath := "./sql_database/race_imports.db"

	err = SeedRaces(dbPath, raceImportDbPath)
	if err != nil {
		t.Errorf("SeedRaces() returned error: %v", err)
	}
}
