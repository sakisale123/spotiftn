package main

import (
<<<<<<< HEAD
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8081"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Users Service! Connected to Mongo at %s", os.Getenv("MONGO_URI"))
	})

	fmt.Printf("Users service starting on port %s\n", port)
=======
	"context"
	"fmt"
	"net/http"
	"os"
	"spotiftn/users/auth"
	"spotiftn/users/handlers"
	"spotiftn/users/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	port := ":8081"

	// Mongo
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	db := client.Database("users")

	// DI
	userRepo := repository.NewUsersRepository(db)
	authService := auth.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Routes
	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)

	fmt.Println("Users service running on", port)
>>>>>>> main
	http.ListenAndServe(port, nil)
}
