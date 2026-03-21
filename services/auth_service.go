package services

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/util"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, loginRequest *models.AuthRequest) (string, *util.APIError) {
	// Retrieve user
	user, err := GetUserByEmail(c, loginRequest.Email)
	if err != nil {
		return "", err
	}

	// Compare password
	passwordValid := CompareUserPassword(user, loginRequest.Password)
	if !passwordValid {
		return "", util.NewAPIError(http.StatusBadRequest, "Invalid username or password")
	}

	// Generate JWT
	jwt, err := GenerateJWTForUser(user)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func IsAdmin(accessToken string) bool {
	claims, err := ParseJWTToken(accessToken)
	if err != nil {
		return false
	}

	roleClaim, ok := claims["role"].(string)
	if !ok {
		return false
	}

	return models.Role(roleClaim) == models.RoleAdmin
}
