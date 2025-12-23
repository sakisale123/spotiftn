package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"spotiftn/content/content_handler"
	"spotiftn/content/repository"
)

func main() {
	ctx := context.Background()

	mongoClient := InitMongoDB(ctx)
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("MongoDB disconnection failed: %v", err)
		}
	}()

	dbName := GetDatabaseName()
	contentRepo := repository.NewMongoContentRepository(mongoClient, dbName)
	contentHandler := content_handler.NewContentHandler(contentRepo)
	router := SetupRoutes(contentHandler)

	port := os.Getenv("SERVER_ADDRESS")

	// ja li tvoje greske treba trazim konju
	if port == "" {
		port = ":8082"
	}

	// Verify if port has : prefix
	addr := port
	if len(port) > 0 && port[0] != ':' {
		addr = ":" + port
	}

	log.Printf("Content Service starting on port %s\n", addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Content Service failed to start: %v", err)
	}
}
