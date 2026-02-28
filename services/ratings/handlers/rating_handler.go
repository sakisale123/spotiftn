package handlers

import (
	"context"
	"net/http"
	"os"
	"time"

	"spotiftn/ratings/client"
	"spotiftn/ratings/db"
	"spotiftn/ratings/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var contentClient *client.ContentClient

func init() {
	contentURL := os.Getenv("CONTENT_SERVICE_URL")
	if contentURL == "" {
		contentURL = "http://localhost:8082"
	}
	contentClient = client.NewContentClient(contentURL)
}

func CreateRating(c *gin.Context) {
	// 1. Get User ID from query or body (to match frontend)
	userID := c.Query("userId")
	if userID == "" {
		userID = "mocked-user-id"
	}

	var req struct {
		SongID string `json:"songId" binding:"required"`
		Score  int    `json:"score" binding:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 2.5 Sinhrona komunikacija između servisa - provera u Content
	exists, err := contentClient.CheckSongExists(req.SongID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify song", "details": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found in Content service"})
		return
	}

	collection := db.Database.Collection("ratings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rating := models.Rating{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		SongID: req.SongID,
		Score:  req.Score,
	}

	_, err = collection.InsertOne(ctx, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rating"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Rating added successfully", "rating": rating})
}

func UpdateRating(c *gin.Context) {
	ratingID := c.Param("id")
	userID := c.Query("userId")
	if userID == "" {
		userID = "mocked-user-id"
	}

	objID, err := primitive.ObjectIDFromHex(ratingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID format"})
		return
	}

	var req struct {
		Score int `json:"score" binding:"required,min=1,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	collection := db.Database.Collection("ratings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objID, "user_id": userID}
	update := bson.M{"$set": bson.M{"score": req.Score}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rating"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found or unathorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating updated successfully"})
}

func DeleteRating(c *gin.Context) {
	ratingID := c.Param("id")
	userID := c.Query("userId")
	if userID == "" {
		userID = "mocked-user-id"
	}

	objID, err := primitive.ObjectIDFromHex(ratingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID format"})
		return
	}

	collection := db.Database.Collection("ratings")
	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer deleteCancel()

	filter := bson.M{"_id": objID, "user_id": userID}

	result, err := collection.DeleteOne(deleteCtx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found or unathorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rating"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rating not found or unathorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rating deleted successfully"})
}

func GetUserRating(c *gin.Context) {
	songID := c.Param("songId")
	userID := c.Query("userId")

	if userID == "" {
		userID = "mocked-user-id"
	}

	collection := db.Database.Collection("ratings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rating models.Rating
	err := collection.FindOne(ctx, bson.M{"song_id": songID, "user_id": userID}).Decode(&rating)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusOK, gin.H{"score": 0})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rating"})
		return
	}

	c.JSON(http.StatusOK, rating)
}
