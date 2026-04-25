package models

type Category struct {
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"not null"`
	UserID uint   `json:"user_id" gorm:"not null;index"`
}
