package services

import (
	"net/http"

	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/storage"
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

	if storage.UserExistsByEmail(email) {
		return util.NewAPIError(http.StatusBadRequest, "User with this email already exists")
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

	if err := storage.CreateUser(&user); err != nil {
		return util.NewAPIError(http.StatusInternalServerError, "Failed to create user.")
	}

	logger.Info("User registration successful.")
	return nil
}

func GetUserByEmail(c *gin.Context, email string) (*models.User, *util.APIError) {
	// Validate Email
	if !util.IsValidEmail(email) {
		return nil, util.NewAPIError(http.StatusBadRequest, "Invalid Email.")
	}

	// Retrieve user based on email
	user, err := storage.FindUserByEmail(email)
	if err != nil {
		return nil, util.NewAPIError(http.StatusBadRequest, err.Error())
	}

	return user, nil
}

// CompareUserPassword returns true if password matches the hashed password
func CompareUserPassword(User *models.User, password string) bool {
	passwordHash := User.Password
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

// GenerateJWTForUser generates a JWT token for a user
func GenerateJWTForUser(user *models.User) (string, *util.APIError) {
	claims := map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	}

	token, err := auth.CreateJWT(claims)
	if err != nil {
		return "", util.NewAPIError(http.StatusInternalServerError, "unable to create access token for user")
	}

	return token, nil
}
