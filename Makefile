.PHONY: drop_db seed_db build run_server

help:
	@echo "drop_db - Remove the database"
	@echo "create_tables - do this first"
	@echo "seed_db - Seed the database"
	@echo "build - Build the project"
	@echo "run_server - Run the server"
	@echo "debug_server - Run the server in dlv"

drop_db:
	rm -f ./sql_database/mud.db

test ?= false

create_tables:
ifeq ($(test), true)
	go run ./internal/seed/create_tables/main.go -test=true
else
	go run ./internal/seed/create_tables/main.go
endif

seed_db:
	go run ./internal/seed/items/main.go
	go run ./internal/seed/areas/main.go
	go run ./internal/seed/display/main.go	
	go run ./internal/seed/players/main.go
	go run ./internal/seed/classes/main.go
	go run ./internal/seed/races/main.go

build:
	go build -o ./bin/mud ./cmd/mud/main.go

run_server:
	go run ./cmd/mud/main.go	


run_character:
		go run ./cmd/character/main.go 


run_ui_list_viewport:
		go run ./cmd/ui_component_testing/list-viewport/main.go

run_ui_tabs_with_list_viewport:
		go run ./cmd/ui_component_testing/tabs/tabs_with_list-viewport/main.go
