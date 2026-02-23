package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"spotiftn/users/auth"
	"spotiftn/users/handlers"
	"spotiftn/users/middleware"
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

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpEmail == "" || smtpPassword == "" {
		fmt.Println("⚠️ WARNING: SMTP configuration is missing or incomplete.")
	} else {
		fmt.Println("✅ SMTP Configuration loaded for:", smtpEmail)
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
	emailService := auth.NewEmailService(smtpHost, smtpPort, smtpEmail, smtpPassword)
	authService := auth.NewAuthService(userRepo, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	// ===== ROUTES =====
	mux := http.NewServeMux()

	// -------- PUBLIC ROUTES --------
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/confirm", authHandler.ConfirmEmail)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/verify-otp", authHandler.VerifyOTP)
	mux.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("/auth/reset-password", authHandler.ResetPassword)

	// -------- PROTECTED ROUTES --------
	mux.Handle(
		"/auth/change-password",
		middleware.AuthMiddleware(http.HandlerFunc(authHandler.ChangePassword)),
	)

	mux.Handle(
		"/auth/logout",
		middleware.AuthMiddleware(http.HandlerFunc(authHandler.Logout)),
	)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("users service OK"))
	})

	// ===== START SERVER WITH RATE LIMIT (DoS PROTECTION) =====
	wrappedMux := middleware.RateLimitMiddleware(mux)

	fmt.Println("Users service running on", port)
	log.Fatal(http.ListenAndServe(port, wrappedMux))
}
