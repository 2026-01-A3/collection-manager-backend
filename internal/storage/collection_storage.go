package storage

import (
	"context"
	"errors"

	"collection-manager-backend/internal/models"

	"gorm.io/gorm"
)

var collectionDB *gorm.DB

func InitCollectionStorage(db *gorm.DB) error {
	collectionDB = db

	return collectionDB.AutoMigrate(&models.Collection{})
}

type BinaryObjectPayload struct {
	Base64    string
	Filename  string
	Extension string
}

func AddCollection(ctx context.Context, name string, categoryID int, bin *BinaryObjectPayload) (models.Collection, error) {
	if collectionDB == nil {
		return models.Collection{}, errors.New("conexão com o banco não inicializada")
	}

	collection := models.Collection{
		Name:       name,
		CategoryID: categoryID,
	}

	if bin != nil {
		saved, err := AddBinaryObject(ctx, bin.Base64, bin.Filename, bin.Extension)
		if err != nil {
			return models.Collection{}, err
		}
		collection.BinaryObjectID = &saved.ID
	}

	if err := collectionDB.WithContext(ctx).Create(&collection).Error; err != nil {
		return models.Collection{}, err
	}

	if err := collectionDB.WithContext(ctx).
		Preload("Category").
		Preload("BinaryObject").
		First(&collection, collection.ID).Error; err != nil {
		return collection, err
	}

	return collection, nil
}

func GetCollections(ctx context.Context) ([]models.Collection, error) {
	if collectionDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	var collections []models.Collection

	if err := collectionDB.WithContext(ctx).
		Preload("Category").
		Preload("BinaryObject").
		Order("id ASC").
		Find(&collections).Error; err != nil {
		return nil, err
	}

	return collections, nil
}

func UpdateCollection(ctx context.Context, id int, name string, categoryID int, bin *BinaryObjectPayload) (models.Collection, error) {
	if collectionDB == nil {
		return models.Collection{}, errors.New("conexão com o banco não inicializada")
	}

	var collection models.Collection
	if err := collectionDB.WithContext(ctx).First(&collection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Collection{}, errors.New("coleção não encontrada")
		}
		return models.Collection{}, err
	}

	collection.Name = name
	collection.CategoryID = categoryID

	if bin != nil {
		if collection.BinaryObjectID != nil {
			updated, err := UpdateBinaryObject(ctx, *collection.BinaryObjectID, bin.Base64, bin.Filename, bin.Extension)
			if err != nil {
				return models.Collection{}, err
			}
			collection.BinaryObjectID = &updated.ID
		} else {
			saved, err := AddBinaryObject(ctx, bin.Base64, bin.Filename, bin.Extension)
			if err != nil {
				return models.Collection{}, err
			}
			collection.BinaryObjectID = &saved.ID
		}
	}

	if err := collectionDB.WithContext(ctx).Save(&collection).Error; err != nil {
		return models.Collection{}, err
	}

	if err := collectionDB.WithContext(ctx).
		Preload("Category").
		Preload("BinaryObject").
		First(&collection, collection.ID).Error; err != nil {
		return collection, err
	}

	return collection, nil
}

func DeleteCollection(ctx context.Context, id int) error {
	if collectionDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var collection models.Collection
	if err := collectionDB.WithContext(ctx).First(&collection, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("coleção não encontrada")
		}
		return err
	}

	binID := collection.BinaryObjectID

	if err := collectionDB.WithContext(ctx).Delete(&collection).Error; err != nil {
		return err
	}

	if binID != nil {
		if err := DeleteBinaryObject(ctx, *binID); err != nil {
			return err
		}
	}

	return nil
}
