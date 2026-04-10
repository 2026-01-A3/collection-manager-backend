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

func AddCategory(ctx context.Context, name string) (models.Category, error) {
	if categoryDB == nil {
		return models.Category{}, errors.New("conexão com o banco não inicializada")
	}

	category := models.Category{
		Name: name,
	}

	if err := categoryDB.WithContext(ctx).Create(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func GetCategories(ctx context.Context) ([]models.Category, error) {
	if categoryDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	var categories []models.Category

	if err := categoryDB.WithContext(ctx).
		Order("id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func UpdateCategory(ctx context.Context, id int, name string) (models.Category, error) {
	if categoryDB == nil {
		return models.Category{}, errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	if err := categoryDB.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Category{}, errors.New("categoria não encontrada")
		}
		return models.Category{}, err
	}

	category.Name = name
	if err := categoryDB.WithContext(ctx).Save(&category).Error; err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func DeleteCategory(ctx context.Context, id int) error {
	if categoryDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var category models.Category
	if err := categoryDB.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("categoria não encontrada")
		}
		return err
	}

	if err := categoryDB.WithContext(ctx).Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
