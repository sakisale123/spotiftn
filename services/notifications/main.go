package main

import (
	"os"
	"spotiftn/notifications/db"
)

func main() {

	db.Init()

	r := SetupRouter()

	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8083"
	}
	r.Run(port)
}
