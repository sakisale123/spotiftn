package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"userId"`         // ID of the user who's subscribing
	TargetID   string             `bson:"target_id" json:"targetId"`     // ID of the Artist or Genre being subscribed to
	TargetType string             `bson:"target_type" json:"targetType"` // "artist" or "genre"
}
