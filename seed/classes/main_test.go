package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestSeedClasses(t *testing.T) {
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

	classImportDbPath := "./sql_database/class_imports.db"

	err = SeedClasses(dbPath, classImportDbPath)
	if err != nil {
		t.Errorf("SeedClasses() returned error: %v", err)
	}
}
