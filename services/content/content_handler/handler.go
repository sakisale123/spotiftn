package content_handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"spotiftn/content/models"
	"spotiftn/content/repository"

	"github.com/colinmarc/hdfs/v2"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	hdfsClient *hdfs.Client
	hdfsOnce   sync.Once
)

type ContentHandler struct {
	Repo repository.ContentRepository
}

func NewContentHandler(repo repository.ContentRepository) *ContentHandler {
	return &ContentHandler{
		Repo: repo,
	}
}

func (h *ContentHandler) getHdfsClient() (*hdfs.Client, error) {
	var err error
	hdfsOnce.Do(func() {
		hdfsAddr := os.Getenv("HDFS_NAMENODE_ADDR")
		if hdfsAddr == "" {
			hdfsAddr = "namenode:9000"
		}
		hdfsClient, err = hdfs.New(hdfsAddr)
		if err == nil {
			log.Println("HDFS client initialized successfully")
		}
	})

	if hdfsClient == nil && err == nil {
		// This happens if the first call failed, we should allow retrying
		hdfsAddr := os.Getenv("HDFS_NAMENODE_ADDR")
		if hdfsAddr == "" {
			hdfsAddr = "namenode:9000"
		}
		hdfsClient, err = hdfs.New(hdfsAddr)
	}

	return hdfsClient, err
}

func (h *ContentHandler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var artist models.Artist
	if err := json.NewDecoder(r.Body).Decode(&artist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if artist.Name == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}

	createdArtist, err := h.Repo.CreateArtist(r.Context(), &artist)
	if err != nil {
		http.Error(w, "Database error: Failed to create artist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdArtist)
}

func (h *ContentHandler) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateArtist models.Artist
	if err := json.NewDecoder(r.Body).Decode(&updateArtist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid Artist ID format", http.StatusBadRequest)
		return
	}
	updateArtist.ID = objID

	if err := h.Repo.UpdateArtist(r.Context(), &updateArtist); err != nil {
		if errors.Is(err, errors.New("artist not found")) {
			http.Error(w, "Artist not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error: Failed to update artist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ContentHandler) GetArtistByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	artist, err := h.Repo.GetArtistByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("artist not found")) {
			http.Error(w, "Artist not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error fetching artist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artist)
}

func (h *ContentHandler) GetAllArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := h.Repo.GetAllArtists(r.Context())
	if err != nil {
		http.Error(w, "Database error: Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

func (h *ContentHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(album.ArtistIDs) == 0 {
		http.Error(w, "Album must be associated with at least one artist", http.StatusBadRequest)
		return
	}

	createdAlbum, err := h.Repo.CreateAlbum(r.Context(), &album)
	if err != nil {
		http.Error(w, "Database error: Failed to create album", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAlbum)
}

func (h *ContentHandler) GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	album, err := h.Repo.GetAlbumByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("album not found")) {
			http.Error(w, "Album not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error fetching album", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func (h *ContentHandler) GetAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistID := vars["id"]

	albums, err := h.Repo.GetAlbumsByArtist(r.Context(), artistID)
	if err != nil {
		http.Error(w, "Database error fetching albums", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

func (h *ContentHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if song.AlbumID.IsZero() {
		http.Error(w, "Song must be associated with an existing album", http.StatusBadRequest)
		return
	}

	if _, err := h.Repo.GetAlbumByID(r.Context(), song.AlbumID.Hex()); err != nil {
		if errors.Is(err, errors.New("album not found")) {
			http.Error(w, "Referenced album does not exist", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error during album check", http.StatusInternalServerError)
		return
	}

	createdSong, err := h.Repo.CreateSong(r.Context(), &song)
	if err != nil {
		http.Error(w, "Database error: Failed to create song", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSong)
}

func (h *ContentHandler) GetSongsByAlbumID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumID := vars["id"]

	songs, err := h.Repo.GetSongsByAlbumID(r.Context(), albumID)
	if err != nil {
		http.Error(w, "Database error fetching songs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h *ContentHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	song, err := h.Repo.GetSongByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errors.New("song not found")) {
			http.Error(w, "Song not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error fetching song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h *ContentHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	types := r.URL.Query()["type"]
	// Clean up empty strings from query parameters
	var cleanTypes []string
	for _, t := range types {
		if t != "" {
			cleanTypes = append(cleanTypes, t)
		}
	}
	types = cleanTypes

	genres := r.URL.Query()["genre"]
	var cleanGenres []string
	for _, g := range genres {
		if g != "" {
			cleanGenres = append(cleanGenres, g)
		}
	}
	genres = cleanGenres

	// If no types specified, search all
	if len(types) == 0 {
		types = []string{"artist", "album", "song"}
	}

	results := make(map[string]interface{})

	for _, t := range types {
		switch t {
		case "artist":
			artists, err := h.Repo.SearchArtists(r.Context(), query, genres)
			if err != nil {
				http.Error(w, "Error searching artists", http.StatusInternalServerError)
				return
			}
			results["artists"] = artists
		case "album":
			albums, err := h.Repo.SearchAlbums(r.Context(), query, genres)
			if err != nil {
				http.Error(w, "Error searching albums", http.StatusInternalServerError)
				return
			}
			results["albums"] = albums
		case "song":
			songs, err := h.Repo.SearchSongs(r.Context(), query, genres)
			if err != nil {
				http.Error(w, "Error searching songs", http.StatusInternalServerError)
				return
			}
			results["songs"] = songs
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *ContentHandler) StreamSong(w http.ResponseWriter, r *http.Request) {
	client, err := h.getHdfsClient()
	if err != nil {
		log.Printf("HDFS client error: %v", err)
		http.Error(w, "HDFS client is not initialized or failed to connect", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	song, err := h.Repo.GetSongByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	if song.AudioURL == "" {
		http.Error(w, "Song has no audio URL", http.StatusNotFound)
		return
	}

	file, err := client.Open(song.AudioURL)
	if err != nil {
		log.Printf("Error opening HDFS file %s: %v", song.AudioURL, err)
		http.Error(w, "Failed to open audio file from HDFS", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/mpeg")
	// Support range requests for seeking
	http.ServeContent(w, r, song.Title+".mp3", time.Time{}, file)
}
