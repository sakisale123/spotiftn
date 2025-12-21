package main

import (
	"log"
	"net/http"
	"spotiftn/content/content_handler"

	"github.com/gorilla/mux"
)

func SetupRoutes(handler *content_handler.ContentHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleware)

	router.HandleFunc("/artists", handler.CreateArtist).Methods("POST")
	router.HandleFunc("/artists/{id}", handler.UpdateArtist).Methods("PUT")
	router.HandleFunc("/albums", handler.CreateAlbum).Methods("POST")
	router.HandleFunc("/songs", handler.CreateSong).Methods("POST")

	router.HandleFunc("/artists", handler.GetAllArtists).Methods("GET")
	router.HandleFunc("/artists/{id}", handler.GetArtistByID).Methods("GET")
	router.HandleFunc("/albums/{id}", handler.GetAlbumByID).Methods("GET")
	router.HandleFunc("/albums", handler.GetAllAlbums).Methods("GET") // [NEW]
	router.HandleFunc("/albums/{id}/songs", handler.GetSongsByAlbumID).Methods("GET")
	router.HandleFunc("/songs/{id}", handler.GetSongByID).Methods("GET") // [NEW]

	router.HandleFunc("/songs/{id}", handler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/albums/{id}", handler.DeleteAlbum).Methods("DELETE")

	return router
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
