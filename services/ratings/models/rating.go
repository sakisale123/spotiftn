package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID string             `bson:"user_id" json:"userId"`
	SongID string             `bson:"song_id" json:"songId"`
	Score  int                `bson:"score" json:"score"` // Assessment score 1-5
}
