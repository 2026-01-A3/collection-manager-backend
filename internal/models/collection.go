package models

type Collection struct {
	ID             int           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string        `json:"name" gorm:"not null"`
	CategoryID     int           `json:"category_id" gorm:"not null"`
	Category       Category      `json:"category" gorm:"foreignKey:CategoryID"`
	BinaryObjectID *int          `json:"binary_object_id"`
	BinaryObject   *BinaryObject `json:"binary_object,omitempty" gorm:"foreignKey:BinaryObjectID"`
}
