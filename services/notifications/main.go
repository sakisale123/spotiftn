package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8083"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Notification Service! Connected to Mongo at %s", os.Getenv("MONGO_URI"))
	})

	fmt.Printf("Notification service starting on port %s\n", port)
	http.ListenAndServe(port, nil)
}
