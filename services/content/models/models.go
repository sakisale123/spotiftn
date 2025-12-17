package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Artist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Biography string             `bson:"biography" json:"biography"`
	Genres    []string           `bson:"genres" json:"genres"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Album struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title     string               `bson:"title" json:"title"`
	Date      time.Time            `bson:"date" json:"date"`
	Genre     string               `bson:"genre" json:"genre"`
	ArtistIDs []primitive.ObjectID `bson:"artist_ids" json:"artist_ids"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}

type Song struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title     string               `bson:"title" json:"title"`
	Duration  int                  `bson:"duration" json:"duration"` // Seconds
	Genre     string               `bson:"genre" json:"genre"`
	AlbumID   primitive.ObjectID   `bson:"album_id" json:"album_id"`
	ArtistIDs []primitive.ObjectID `bson:"artist_ids" json:"artist_ids"`
	AudioURL  string               `bson:"audio_url" json:"audio_url"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
}
