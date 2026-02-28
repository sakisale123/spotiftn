package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"spotiftn/subscriptions/client"
	"spotiftn/subscriptions/db"
	"spotiftn/subscriptions/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var contentClient *client.ContentClient

func init() {
	contentURL := os.Getenv("CONTENT_SERVICE_URL")
	if contentURL == "" {
		contentURL = "https://content:8082" // Default internal address for content
	}
	contentClient = client.NewContentClient(contentURL)
}

func Subscribe(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = "mocked-user-id"
	}

	var req struct {
		TargetID   string `json:"targetId" binding:"required"`
		TargetType string `json:"targetType" binding:"required,oneof=artist genre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 2.5 Sinhrona komunikacija između servisa - provera u Content
	if req.TargetType == "artist" {
		exists, err := contentClient.CheckArtistExists(req.TargetID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify artist existence", "details": err.Error()})
			return
		}

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found in Content service"})
			return
		}
	}

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
	userID := c.Query("userId")
	if userID == "" {
		userID = "mocked-user-id"
	}
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

func GetSubscriptionStatus(c *gin.Context) {
	targetID := c.Query("targetId")
	userID := c.Query("userId")

	if userID == "" {
		userID = "mocked-user-id"
	}

	collection := db.Database.Collection("subscription")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sub models.Subscription
	err := collection.FindOne(ctx, bson.M{"user_id": userID, "target_id": targetID}).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"isSubscribed": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isSubscribed": true, "subscriptionId": sub.ID.Hex()})
}
