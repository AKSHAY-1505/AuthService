package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/AKSHAY-1505/auth-service/auth"
	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/AKSHAY-1505/auth-service/models"
	"github.com/AKSHAY-1505/auth-service/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	if userAlreadyExists(email) {
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

	result := initializers.DB.Create(&user)
	if result.Error != nil {
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
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// user not found
			return nil, util.NewAPIError(http.StatusBadRequest, "User is not registered.")
		}

		return nil, util.NewAPIError(http.StatusInternalServerError, "Unable to find user.")
	}

	return &user, nil
}

// ComparePassword returns true if password matches the hashed password
func CompareUserPassword(User *models.User, password string) bool {
	passwordHash := User.Password
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token for a user
func GenerateJWTForUser(user *models.User) (string, *util.APIError) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // token expires in 24h
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(auth.GetPrivateKey())
	if err != nil {
		return "", util.NewAPIError(http.StatusInternalServerError, err.Error())
	}

	return signedToken, nil
}

func ParseJWTToken(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Enforce expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.GetPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Optional: manually verify exp if you want stricter control
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func userAlreadyExists(email string) bool {
	var exists bool

	err := initializers.DB.
		Model(&models.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).Error

	if err != nil {
		return false
	}

	return exists
}
