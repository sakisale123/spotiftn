package content_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"spotiftn/content/events"
	"spotiftn/content/models"
	"spotiftn/content/repository"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContentHandler struct {
	Repo      repository.ContentRepository
	Publisher events.EventPublisher
}

func NewContentHandler(repo repository.ContentRepository, publisher events.EventPublisher) *ContentHandler {
	return &ContentHandler{
		Repo:      repo,
		Publisher: publisher,
	}
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

	// Publish Event
	if h.Publisher != nil {
		event := events.ArtistCreatedEvent{
			ID:        createdArtist.ID.Hex(),
			Name:      createdArtist.Name,
			Genres:    createdArtist.Genres,
			Timestamp: time.Now(),
		}
		h.Publisher.Publish(events.SubjectArtistCreated, event)
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

	// Publish Event
	if h.Publisher != nil {
		event := events.ArtistUpdatedEvent{
			ID:        id,
			Name:      updateArtist.Name,
			Genres:    updateArtist.Genres,
			Timestamp: time.Now(),
		}
		h.Publisher.Publish(events.SubjectArtistUpdated, event)
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

	if album.Title == "" {
		http.Error(w, "Album title is required", http.StatusBadRequest)
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

	// Publish Event
	if h.Publisher != nil {
		artistIDs := make([]string, len(createdAlbum.ArtistIDs))
		for i, aid := range createdAlbum.ArtistIDs {
			artistIDs[i] = aid.Hex()
		}
		event := events.AlbumCreatedEvent{
			ID:        createdAlbum.ID.Hex(),
			Title:     createdAlbum.Title,
			ArtistIDs: artistIDs,
			Timestamp: time.Now(),
		}
		h.Publisher.Publish(events.SubjectAlbumCreated, event)
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

func (h *ContentHandler) GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := h.Repo.GetAllAlbums(r.Context())
	if err != nil {
		http.Error(w, "Database error: Failed to fetch albums", http.StatusInternalServerError)
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

	if song.Title == "" {
		http.Error(w, "Song title is required", http.StatusBadRequest)
		return
	}
	if song.Duration <= 0 {
		http.Error(w, "Song duration must be positive", http.StatusBadRequest)
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

	// Publish Event
	if h.Publisher != nil {
		artistIDs := make([]string, len(createdSong.ArtistIDs))
		for i, aid := range createdSong.ArtistIDs {
			artistIDs[i] = aid.Hex()
		}
		event := events.SongCreatedEvent{
			ID:        createdSong.ID.Hex(),
			Title:     createdSong.Title,
			AlbumID:   createdSong.AlbumID.Hex(),
			ArtistIDs: artistIDs,
			Timestamp: time.Now(),
		}
		h.Publisher.Publish(events.SubjectSongCreated, event)
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
