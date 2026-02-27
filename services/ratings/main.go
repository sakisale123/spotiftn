package main

import (
	"os"
	"spotiftn/ratings/db"
)

func main() {
	db.Init()

	r := SetupRouter()

	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8087"
	}
	// RunTLS for HTTPS
	r.RunTLS(port, "/etc/ssl/certs/server.crt", "/etc/ssl/certs/server.key")
}
