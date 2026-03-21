package handlers

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/services"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var registerRequest models.AuthRequest

	err := c.ShouldBindBodyWithJSON(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body.",
		})
		return
	}

	// Create the user
	createError := services.CreateUser(c, registerRequest.Email, registerRequest.Password, models.RoleUser)
	if createError != nil {
		c.JSON(createError.StatusCode, gin.H{
			"error": createError.Message,
		})
		return
	}

	// Send response
	c.JSON(http.StatusCreated, gin.H{})
}
