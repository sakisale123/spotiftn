package main

import (
	"log"
	"net/http"
	"spotiftn/content/content_handler"
	"spotiftn/content/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(handler *content_handler.ContentHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	// Protected routes
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/artists", handler.CreateArtist).Methods("POST")
	protected.HandleFunc("/artists/{id}", handler.UpdateArtist).Methods("PUT")
	protected.HandleFunc("/albums", handler.CreateAlbum).Methods("POST")
	protected.HandleFunc("/songs", handler.CreateSong).Methods("POST")

	router.HandleFunc("/artists", handler.GetAllArtists).Methods("GET")
	router.HandleFunc("/artists/{id}", handler.GetArtistByID).Methods("GET")
	router.HandleFunc("/artists/{id}/albums", handler.GetAlbumsByArtist).Methods("GET")
	router.HandleFunc("/albums/{id}", handler.GetAlbumByID).Methods("GET")
	router.HandleFunc("/albums/{id}/songs", handler.GetSongsByAlbumID).Methods("GET")

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
