package events

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	SubjectArtistCreated = "artist.created"
	SubjectArtistUpdated = "artist.updated"
	SubjectAlbumCreated  = "album.created"
	SubjectSongCreated   = "song.created"
)

type ArtistCreatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Genres    []string  `json:"genres"`
	Timestamp time.Time `json:"timestamp"`
}

type ArtistUpdatedEvent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Genres    []string  `json:"genres"`
	Timestamp time.Time `json:"timestamp"`
}

type AlbumCreatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	ArtistIDs []string  `json:"artist_ids"`
	Timestamp time.Time `json:"timestamp"`
}

type SongCreatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	AlbumID   string    `json:"album_id"`
	ArtistIDs []string  `json:"artist_ids"`
	Timestamp time.Time `json:"timestamp"`
}

type EventPublisher interface {
	Publish(subject string, data interface{}) error
	Close()
}
type NatsEventPublisher struct {
	conn *nats.Conn
}

func NewNatsEventPublisher(url string) (*NatsEventPublisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventPublisher{conn: conn}, nil
}

func (p *NatsEventPublisher) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.conn.Publish(subject, payload)
}

func (p *NatsEventPublisher) Close() {
	if p.conn != nil {
		p.conn.Close()
	}
}
