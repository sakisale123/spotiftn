package main

import (
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
	http.ListenAndServe(port, nil)
}
