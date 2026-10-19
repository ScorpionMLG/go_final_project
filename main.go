package main

import (
	"log"
	"main/pkg/db"
	"main/pkg/server"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Critical error: Failed to initialize database: %v", err)
	}
	defer db.Close()
	err = server.Run()
	if err != nil {
		log.Fatalf("Critical error: Failed to initialize server: %v", err)
	}
}
