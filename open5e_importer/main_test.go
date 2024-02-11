package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestImportMonsters(t *testing.T) {
	err := os.Remove("../sql_database/test.db")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", "../sql_database/test.db")
	if err != nil {
		log.Fatalf("Failed to open sqlite db: %v", err)
	}
	if err != nil {
		if os.IsNotExist(err) {
			// don't worry about it!
		}
		t.Fatal(err)
	}
	data, err := ioutil.ReadFile("./test_data/testdata.json")
	if err != nil {
		t.Fatal(err)
	}
	monsters, _ := convertJsonToMonsterImports(data)
	writeMonstersToDB(db, monsters)
}
