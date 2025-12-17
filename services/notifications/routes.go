package main

import (
	"spotiftn/notifications/handlers"
<<<<<<< HEAD
=======
	"spotiftn/notifications/middleware"
>>>>>>> main

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

<<<<<<< HEAD
	r.GET("/:userID", handlers.GetNotifications)
	r.POST("/", handlers.CreateNotification)

=======
>>>>>>> main
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "notifications"})
	})

<<<<<<< HEAD
=======
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/:userID", handlers.GetNotifications)
		api.POST("/", handlers.CreateNotification)
	}

>>>>>>> main
	return r
}
