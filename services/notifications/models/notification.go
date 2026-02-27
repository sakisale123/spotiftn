package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Notification struct {
	ID        gocql.UUID `json:"id"`
	UserID    string     `json:"user_id" binding:"required"`
	Type      string     `json:"type" binding:"required"`
	Message   string     `json:"message" binding:"required"`
	CreatedAt time.Time  `json:"created_at"`
	IsRead    bool       `json:"is_read"`
}
