.PHONY: drop_db seed_db build run_server

help:
	@echo "drop_db - Remove the database"
	@echo "create_tables - do this first"
	@echo "seed_db - Seed the database"
	@echo "build - Build the project"
	@echo "run_server - Run the server"
	@echo "debug_server - Run the server in dlv"
	@echo "examine_actions - debugging stuff"
	@echo "import_races"
	@echo "import_classes"
	@echo "import_monsters"

drop_db:
	rm -f ./sql_database/mud.db


test ?= false

create_tables:
ifeq ($(test), true)
	go run ./seed/create_tables/main.go -test=true
else
	go run ./seed/create_tables/main.go
endif

seed_db:
	go run ./seed/items/main.go
	go run ./seed/areas/main.go
	go run ./seed/display/main.go	
	go run ./seed/players/main.go
	go run ./seed/classes/main.go
	go run ./seed/races/main.go

build:
	go build -o ./bin/mud ./main.go

run_server:
	go run ./main.go	

debug_server:
	dlv debug -l 127.0.0.1:38697 --headless ./main.go
