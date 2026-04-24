package storage

import (
	"context"
	"errors"

	"collection-manager-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var userDB *gorm.DB

func InitUserStorage(db *gorm.DB) error {
	userDB = db
	return userDB.AutoMigrate(&models.User{})
}

func CreateUser(ctx context.Context, name, email, password string, role models.Role) (models.User, error) {
	if userDB == nil {
		return models.User{}, errors.New("conexão com o banco não inicializada")
	}

	if role == "" {
		role = models.UserRole
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := userDB.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if userDB == nil {
		return models.User{}, errors.New("conexão com o banco não inicializada")
	}

	var user models.User
	if err := userDB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
