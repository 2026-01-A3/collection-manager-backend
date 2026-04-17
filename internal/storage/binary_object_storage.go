package storage

import (
	"context"
	"errors"

	"collection-manager-backend/internal/models"

	"gorm.io/gorm"
)

var binaryObjectDB *gorm.DB

func InitBinaryObjectStorage(db *gorm.DB) error {
	binaryObjectDB = db

	return binaryObjectDB.AutoMigrate(&models.BinaryObject{})
}

func AddBinaryObject(ctx context.Context, base64, filename, extension string) (models.BinaryObject, error) {
	if binaryObjectDB == nil {
		return models.BinaryObject{}, errors.New("conexão com o banco não inicializada")
	}

	bin := models.BinaryObject{
		Base64:    base64,
		Filename:  filename,
		Extension: extension,
	}

	if err := binaryObjectDB.WithContext(ctx).Create(&bin).Error; err != nil {
		return models.BinaryObject{}, err
	}

	return bin, nil
}

func UpdateBinaryObject(ctx context.Context, id int, base64, filename, extension string) (models.BinaryObject, error) {
	if binaryObjectDB == nil {
		return models.BinaryObject{}, errors.New("conexão com o banco não inicializada")
	}

	var bin models.BinaryObject
	if err := binaryObjectDB.WithContext(ctx).First(&bin, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.BinaryObject{}, errors.New("arquivo não encontrado")
		}
		return models.BinaryObject{}, err
	}

	bin.Base64 = base64
	bin.Filename = filename
	bin.Extension = extension

	if err := binaryObjectDB.WithContext(ctx).Save(&bin).Error; err != nil {
		return models.BinaryObject{}, err
	}

	return bin, nil
}

func DeleteBinaryObject(ctx context.Context, id int) error {
	if binaryObjectDB == nil {
		return errors.New("conexão com o banco não inicializada")
	}

	return binaryObjectDB.WithContext(ctx).Delete(&models.BinaryObject{}, id).Error
}
