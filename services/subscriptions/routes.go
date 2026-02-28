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

	r.POST("/", handlers.Subscribe)
	r.DELETE("/:id", handlers.Unsubscribe)
	r.GET("/status", handlers.GetSubscriptionStatus)

	return r
}
