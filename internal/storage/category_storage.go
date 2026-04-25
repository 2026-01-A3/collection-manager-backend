package storage

import (
	"context"
	"errors"

	"collection-manager-backend/internal/models"

	"gorm.io/gorm"
)

var categoryDB *gorm.DB

func InitCategoryStorage(db *gorm.DB) error {
	categoryDB = db

	return categoryDB.AutoMigrate(&models.Category{})
}

func AddCategory(ctx context.Context, userID uint, name string) (models.Category, error) {
	if categoryDB == nil {
		return models.Category{}, errors.New("conexão com o banco não inicializada")
	}

	category := models.Category{
		Name:   name,
		UserID: userID,
	}

	if err := categoryDB.WithContext(ctx).Create(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func GetCategories(ctx context.Context, userID uint) ([]models.Category, error) {
	if categoryDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	var categories []models.Category

	if err := categoryDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func UpdateCategory(ctx context.Context, userID uint, id int, name string) (models.Category, error) {
	if categoryDB == nil {
		return models.Category{}, errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	if err := categoryDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, ErrNotFound
		}
		return models.Category{}, err
	}

	category.Name = name
	if err := categoryDB.WithContext(ctx).Save(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func DeleteCategory(ctx context.Context, userID uint, id int) error {
	if categoryDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	if err := categoryDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	if err := categoryDB.WithContext(ctx).Delete(&category).Error; err != nil {
		return err
	}

	return nil
}

// CategoryBelongsToUser verifica se uma categoria pertence ao usuário.
// Usado para validar FK ao criar/atualizar coleções.
func CategoryBelongsToUser(ctx context.Context, userID uint, categoryID int) (bool, error) {
	if categoryDB == nil {
		return false, errors.New("conexão com o banco não inicializada")
	}

	var count int64
	if err := categoryDB.WithContext(ctx).
		Model(&models.Category{}).
		Where("id = ? AND user_id = ?", categoryID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
