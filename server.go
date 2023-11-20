package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"mud/areas"
	"mud/commands"
	"mud/interfaces"
	"mud/player"
	"net"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type CommandRouterInterface interface {
	HandleCommand(db *sql.DB, player interfaces.PlayerInterface, command []byte)
}

func handleConnection(conn net.Conn, router CommandRouterInterface, db *sql.DB) {
	defer conn.Close()

	player := player.NewPlayer(conn)

	if player.GetName() == "" {
		fmt.Fprintf(conn, "Welcome! Please enter your player name: ")
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		player.Name = strings.TrimSpace(string(buf[:n]))
	}

	// Retrieve player info from the database
	err := db.QueryRow("SELECT uuid, area, room, health FROM players WHERE name = ?", player.Name).
		Scan(&player.UUID, &player.Area, &player.Room, &player.Health)
	if err != nil {
		fmt.Fprintf(conn, "Error retrieving player info: %v\n", err)
		return
	}

	if player.Area == "" || player.Room == "" {
		player.Area = "d71e8cf1-d5ba-426c-8915-4c7f5b22e3a9"
		player.Room = "189a729d-4e40-4184-a732-e2c45c66ff46"
	}

	router.HandleCommand(db, player, bytes.NewBufferString("look").Bytes())

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		router.HandleCommand(db, player, buf[:n])
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./sql_database/mud.db")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	} else {
		err := db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Database opened successfully")
	}
	player.SeedPlayers(db)
	areas.SeedAreasAndRooms(db)

	defer db.Close()

	router := commands.NewCommandRouter()

	commands.RegisterCommands(router, commands.CommandHandlers)

	// Check if the command router is empty.
	if len(router.Handlers) == 0 {
		fmt.Println("Warning: no commands registered. Exiting...")
		return
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()

	wg := sync.WaitGroup{}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		wg.Add(1)

		go handleConnection(conn, router, db)
	}

	wg.Wait()
}
