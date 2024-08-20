package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestSeedAreasAndRooms(t *testing.T) {
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
	cmd := exec.Command("make", "create_mobs_table", "test=true")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("some kinda error? %v", err)
	}

	monstersImportDbPath := "./sql_database/monster_imports.db"

	err = SeedAreasAndRooms(dbPath, monstersImportDbPath)
	if err != nil {
		t.Errorf("SeedAreasAndRooms() returned error: %v", err)
	}
}
