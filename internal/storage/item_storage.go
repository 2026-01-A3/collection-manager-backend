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

func AddItem(ctx context.Context, name string, price float64, collectionID int, bin *BinaryObjectPayload) (models.Item, error) {
	if itemDB == nil {
		return models.Item{}, errors.New("conexão com o banco não inicializada")
	}

	item := models.Item{
		Name:         name,
		Price:        price,
		CollectionID: collectionID,
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

func GetItemsByCollection(ctx context.Context, collectionID int) ([]models.Item, error) {
	if itemDB == nil {
		return nil, errors.New("conexão com o banco não inicializada")
	}

	var items []models.Item

	if err := itemDB.WithContext(ctx).
		Preload("BinaryObject").
		Where("collection_id = ?", collectionID).
		Order("id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func UpdateItem(ctx context.Context, id int, name string, price float64, collectionID int, bin *BinaryObjectPayload) (models.Item, error) {
	if itemDB == nil {
		return models.Item{}, errors.New("conexão com o banco não inicializada")
	}

	var item models.Item
	if err := itemDB.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Item{}, errors.New("item não encontrado")
		}
		return models.Item{}, err
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

func DeleteItem(ctx context.Context, id int) error {
	if itemDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	var item models.Item
	if err := itemDB.WithContext(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("item não encontrado")
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
