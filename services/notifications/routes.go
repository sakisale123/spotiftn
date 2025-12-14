package main

import (
	"spotiftn/notifications/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/:userID", handlers.GetNotifications)
	r.POST("/", handlers.CreateNotification)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "notifications"})
	})

	return r
}
