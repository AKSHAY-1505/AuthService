package controllers

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/services"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	unmarshalErr := c.ShouldBindBodyWithJSON(&loginRequest)
	if unmarshalErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body.",
		})
		return
	}

	// Retrieve user
	user, err := services.GetUserByEmail(c, loginRequest.Email)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"error": err.Message,
		})

		return
	}

	// Compare password
	passwordValid := services.CompareUserPassword(user, loginRequest.Password)
	if !passwordValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password.",
		})

		return
	}

	// Generate JWT
	jwt, err := services.GenerateJWTForUser(user)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"error": err.Message,
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"token": jwt,
	})
}
