package handlers

import (
	"context"
	"net/http"
	"time"

	"spotiftn/subscriptions/db"
	"spotiftn/subscriptions/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Subscribe(c *gin.Context) {
	// userID := c.GetString("userId")
	userID := "mocked-user-id"

	var req struct {
		TargetID   string `json:"targetId" binding:"required"`
		TargetType string `json:"targetType" binding:"required,oneof=artist genre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// TODO: Member 3 Check - Sync HTTP call to check if artist or genre exists

	collection := db.Database.Collection("subscription")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Proveravamo da li vec postoji pretplata
	filter := bson.M{
		"user_id":     userID,
		"target_id":   req.TargetID,
		"target_type": req.TargetType,
	}

	var existing models.Subscription
	err := collection.FindOne(ctx, filter).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already subscribed to this " + req.TargetType})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking subscription"})
		return
	}

	subscription := models.Subscription{
		ID:         primitive.NewObjectID(),
		UserID:     userID,
		TargetID:   req.TargetID,
		TargetType: req.TargetType,
	}

	_, err = collection.InsertOne(ctx, subscription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subscribed successfully", "subscription": subscription})
}

func Unsubscribe(c *gin.Context) {
	// userID := c.GetString("userId")
	userID := "mocked-user-id"
	targetID := c.Param("id")

	collection := db.Database.Collection("subscription")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":   userID,
		"target_id": targetID,
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
