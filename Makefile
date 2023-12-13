.PHONY: drop_db seed_db build run_server

help:
	@echo "drop_db - Remove the database"
	@echo "seed_db - Seed the database"
	@echo "build - Build the project"
	@echo "run_server - Run the server"

drop_db:
	rm -f ./sql_database/mud.db

seed_db:
	go run ./seed/areas/main.go
	go run ./seed/display/main.go	
	go run ./seed/players/main.go

build:
	go build -o ./bin/mud ./main.go

run_server:
	go run ./main.go	
