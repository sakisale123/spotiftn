package handlers

import (
	"net/http"
	"sort"
	"time"

	"spotiftn/notifications/db"
	"spotiftn/notifications/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func GetNotifications(c *gin.Context) {
	userID := c.Param("userID")

	notifications := make([]models.Notification, 0)

	fetchForUser := func(uid string) {
		iter := db.Session.Query(`SELECT id, user_id, type, message, created_at, is_read FROM notifications WHERE user_id = ?`, uid).Iter()
		var notif models.Notification
		for iter.Scan(&notif.ID, &notif.UserID, &notif.Type, &notif.Message, &notif.CreatedAt, &notif.IsRead) {
			notifications = append(notifications, notif)
		}
		if err := iter.Close(); err != nil {

		}
	}

	fetchForUser(userID)
	fetchForUser("global")

	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].CreatedAt.After(notifications[j].CreatedAt)
	})

	c.JSON(http.StatusOK, notifications)
}

func CreateNotification(c *gin.Context) {
	var notif models.Notification
	if err := c.ShouldBindJSON(&notif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notif.ID, _ = gocql.RandomUUID()
	notif.CreatedAt = time.Now()
	notif.IsRead = false

	err := db.Session.Query(`INSERT INTO notifications (id, user_id, type, message, created_at, is_read) VALUES (?, ?, ?, ?, ?, ?)`,
		notif.ID, notif.UserID, notif.Type, notif.Message, notif.CreatedAt, notif.IsRead).Exec()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, notif)
}
