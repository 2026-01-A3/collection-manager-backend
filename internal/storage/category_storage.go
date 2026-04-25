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

func scopeByUser(tx *gorm.DB, userID uint, isAdmin bool) *gorm.DB {
	if isAdmin {
		return tx
	}
	return tx.Where("user_id = ?", userID)
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

func GetCategories(ctx context.Context, userID uint, isAdmin bool) ([]models.Category, error) {
	if categoryDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	var categories []models.Category

	tx := scopeByUser(categoryDB.WithContext(ctx), userID, isAdmin)
	if err := tx.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func UpdateCategory(ctx context.Context, userID uint, isAdmin bool, id int, name string) (models.Category, error) {
	if categoryDB == nil {
		return models.Category{}, errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	tx := scopeByUser(categoryDB.WithContext(ctx), userID, isAdmin)
	if err := tx.Where("id = ?", id).First(&category).Error; err != nil {
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

func DeleteCategory(ctx context.Context, userID uint, isAdmin bool, id int) error {
	if categoryDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	tx := scopeByUser(categoryDB.WithContext(ctx), userID, isAdmin)
	if err := tx.Where("id = ?", id).First(&category).Error; err != nil {
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

func CategoryAccessible(ctx context.Context, userID uint, isAdmin bool, categoryID int) (bool, error) {
	if categoryDB == nil {
		return false, errors.New("conexão com o banco não inicializada")
	}

	var count int64
	tx := scopeByUser(categoryDB.WithContext(ctx).Model(&models.Category{}), userID, isAdmin)
	if err := tx.Where("id = ?", categoryID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
