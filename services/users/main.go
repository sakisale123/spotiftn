package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"spotiftn/users/auth"
	"spotiftn/users/handlers"
	"spotiftn/users/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// ===== CONFIG =====
	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = ":8081"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "users_db"
	}

	// ===== MONGO =====
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	// ===== DEPENDENCY INJECTION =====
	userRepo := repository.NewUsersRepository(db)
	authService := auth.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// ===== ROUTES =====
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/confirm", authHandler.ConfirmEmail)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/verify-otp", authHandler.VerifyOTP)

	mux.HandleFunc("/auth/logout", authHandler.Logout)

	// Health check (korisno za test)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("users service OK"))
	})

	// ===== START SERVER =====
	fmt.Println("Users service running on", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
