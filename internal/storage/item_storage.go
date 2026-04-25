package storage

import (
	"context"
	"errors"

	"collection-manager-backend/internal/models"

	"gorm.io/gorm"
)

var itemDB *gorm.DB

func InitItemStorage(db *gorm.DB) error {
	itemDB = db

	return itemDB.AutoMigrate(&models.Item{})
}

func AddItem(ctx context.Context, userID uint, name string, price float64, collectionID int, bin *BinaryObjectPayload) (models.Item, error) {
	if itemDB == nil {
		return models.Item{}, errors.New("conexão com o banco não inicializada")
	}

	owns, err := CollectionBelongsToUser(ctx, userID, collectionID)
	if err != nil {
		return models.Item{}, err
	}
	if !owns {
		return models.Item{}, ErrNotFound
	}

	item := models.Item{
		Name:         name,
		Price:        price,
		CollectionID: collectionID,
		UserID:       userID,
	}

	if bin != nil {
		saved, err := AddBinaryObject(ctx, bin.Base64, bin.Filename, bin.Extension)
		if err != nil {
			return models.Item{}, err
		}
		item.BinaryObjectID = &saved.ID
	}

	if err := itemDB.WithContext(ctx).Create(&item).Error; err != nil {
		return models.Item{}, err
	}

	if err := itemDB.WithContext(ctx).
		Preload("BinaryObject").
		First(&item, item.ID).Error; err != nil {
		return item, err
	}

	return item, nil
}

func GetItemsByCollection(ctx context.Context, userID uint, collectionID int) ([]models.Item, error) {
	if itemDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	owns, err := CollectionBelongsToUser(ctx, userID, collectionID)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, ErrNotFound
	}

	var items []models.Item

	if err := itemDB.WithContext(ctx).
		Preload("BinaryObject").
		Where("collection_id = ? AND user_id = ?", collectionID, userID).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func UpdateItem(ctx context.Context, userID uint, id int, name string, price float64, collectionID int, bin *BinaryObjectPayload) (models.Item, error) {
	if itemDB == nil {
		return models.Item{}, errors.New("conexão com o banco não inicializada")
	}

	var item models.Item
	if err := itemDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Item{}, ErrNotFound
		}
		return models.Item{}, err
	}

	owns, err := CollectionBelongsToUser(ctx, userID, collectionID)
	if err != nil {
		return models.Item{}, err
	}
	if !owns {
		return models.Item{}, ErrNotFound
	}

	item.Name = name
	item.Price = price
	item.CollectionID = collectionID

	if bin != nil {
		if item.BinaryObjectID != nil {
			updated, err := UpdateBinaryObject(ctx, *item.BinaryObjectID, bin.Base64, bin.Filename, bin.Extension)
			if err != nil {
				return models.Item{}, err
			}
			item.BinaryObjectID = &updated.ID
		} else {
			saved, err := AddBinaryObject(ctx, bin.Base64, bin.Filename, bin.Extension)
			if err != nil {
				return models.Item{}, err
			}
			item.BinaryObjectID = &saved.ID
		}
	}

	if err := itemDB.WithContext(ctx).Save(&item).Error; err != nil {
		return models.Item{}, err
	}

	if err := itemDB.WithContext(ctx).
		Preload("BinaryObject").
		First(&item, item.ID).Error; err != nil {
		return item, err
	}

	return item, nil
}

func DeleteItem(ctx context.Context, userID uint, id int) error {
	if itemDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var item models.Item
	if err := itemDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	binID := item.BinaryObjectID

	if err := itemDB.WithContext(ctx).Delete(&item).Error; err != nil {
		return err
	}

	if binID != nil {
		if err := DeleteBinaryObject(ctx, *binID); err != nil {
			return err
		}
	}

	return nil
}
