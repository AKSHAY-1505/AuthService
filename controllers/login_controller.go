package controllers

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/services"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginRequest models.AuthRequest

	unmarshalErr := c.ShouldBindBodyWithJSON(&loginRequest)
	if unmarshalErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body.",
		})
		return
	}

	jwt, err := services.Login(c, &loginRequest)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"error": err.Message,
		})
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
	})
}
