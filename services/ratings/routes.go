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

	r.POST("/", handlers.CreateRating)
	r.PUT("/:id", handlers.UpdateRating)
	r.DELETE("/:id", handlers.DeleteRating)
	r.GET("/user/song/:songId", handlers.GetUserRating)

	return r
}
