package handlers

import (
	"context"
	"net/http"
	"time"

	"spotiftn/notifications/db"
	"spotiftn/notifications/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNotifications(c *gin.Context) {
	userID := c.Param("userID")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := db.Collection.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer cursor.Close(ctx)

<<<<<<< HEAD
	var notifications []models.Notification
=======
	notifications := make([]models.Notification, 0)

>>>>>>> main
	if err = cursor.All(ctx, &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Parsing error"})
		return
	}

<<<<<<< HEAD
	if notifications == nil {
		notifications = []models.Notification{}
	}

=======
>>>>>>> main
	c.JSON(http.StatusOK, notifications)
}

func CreateNotification(c *gin.Context) {
	var notif models.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif.ID = primitive.NewObjectID()
	notif.CreatedAt = time.Now()
	notif.IsRead = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Collection.InsertOne(ctx, notif)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}

	c.JSON(http.StatusCreated, notif)
}
