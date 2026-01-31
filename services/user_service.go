package services

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context, email string, password string, role models.Role) *util.APIError {
	logger := util.GetLoggerFromContext(c)

	// validate Email & Password
	if !util.IsValidEmail(email) {
		return util.NewAPIError(http.StatusBadRequest, "Invalid Email Format.")
	}

	if password == "" {
		return util.NewAPIError(http.StatusBadRequest, "Password cannot be empty.")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return util.NewAPIError(http.StatusInternalServerError, "Unable to create user.")
	}

	// Create User
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		return util.NewAPIError(http.StatusInternalServerError, "Failed to create user.")
	}

	logger.Info("User registration successful.")
	return nil
}
