package main

import (
	"spotiftn/subscriptions/handlers"

	"github.com/gin-gonic/gin"
	// "spotiftn/subscriptions/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "subscriptions"})
	})

	api := r.Group("/")
	// api.Use(middleware.AuthMiddleware()) // TODO test auth middleware later
	{
		api.POST("/", handlers.Subscribe)
		api.DELETE("/:id", handlers.Unsubscribe)
	}

	return r
}
