package storage

import (
	"errors"

	"github.com/AKSHAY-1505/auth-service/initializers"
	"github.com/AKSHAY-1505/auth-service/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	result := initializers.DB.Create(user)
	return result.Error
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := initializers.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// user not found
			return nil, errors.New("user not found")
		}

		return nil, errors.New("unable to find user")
	}

	return &user, nil
}

func UserExistsByEmail(email string) bool {
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
