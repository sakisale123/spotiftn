package main

import (
	"os"
	"spotiftn/subscriptions/db"
)

func main() {
	db.Init()

	r := SetupRouter()

	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8089"
	}
	// RunTLS for HTTPS
	r.RunTLS(port, "/etc/ssl/certs/server.crt", "/etc/ssl/certs/server.key")
}
