package services

import (
	"api-go-test/models"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (us *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User

	if err := us.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	// Remove passwords from all users for security
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (us *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	if err := us.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Remove password from response for security
	user.Password = ""

	return &user, nil
}
