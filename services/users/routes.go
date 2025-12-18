package main

import (
	"net/http"

	"spotiftn/users/handlers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
) {
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
}
