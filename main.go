package main

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.InitLogger()
	initializers.InitEnvVariables()
	initializers.InitDB()
}

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
