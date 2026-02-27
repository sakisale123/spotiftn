package main

import (
	"spotiftn/ratings/handlers"

	"github.com/gin-gonic/gin"
	// "spotiftn/ratings/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "ratings"})
	})

	api := r.Group("/")
	// api.Use(middleware.AuthMiddleware()) // TODO: implement and test Auth middleware later
	{
		api.POST("/", handlers.CreateRating)
		api.PUT("/:id", handlers.UpdateRating)
		api.DELETE("/:id", handlers.DeleteRating)
	}

	return r
}
