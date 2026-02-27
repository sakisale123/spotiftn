package main

import (
	"spotiftn/notifications/handlers"
	"spotiftn/notifications/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "notifications"})
	})

	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/:userID", handlers.GetNotifications)
		api.POST("/", handlers.CreateNotification)
	}

	return r
}
