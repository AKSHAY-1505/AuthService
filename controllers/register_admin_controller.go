package controllers

import (
	"net/http"
	"strings"

	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/services"
	"github.com/gin-gonic/gin"
)

func RegisterAdmin(c *gin.Context) {
	// 1. Extract and Validate Token
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := services.ParseJWTToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid jwt token"})
		return
	}

	roleClaim, ok := claims["role"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role in token"})
		return
	}

	if models.Role(roleClaim) != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorised"})
		return
	}

	// Check for valid token with ADMIN role

	var registerRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = c.ShouldBindBodyWithJSON(&registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body.",
		})
		return
	}

	// Create the user
	createError := services.CreateUser(c, registerRequest.Email, registerRequest.Password, models.RoleAdmin)
	if createError != nil {
		c.JSON(createError.StatusCode, gin.H{
			"error": createError.Message,
		})
		return
	}

	// Send response
	c.JSON(http.StatusCreated, gin.H{})
}
