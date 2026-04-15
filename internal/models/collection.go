package models

type Collection struct {
	ID         int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"not null"`
	CategoryID int      `json:"category_id" gorm:"not null"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	ImageURL   string   `json:"image_url"`
}
