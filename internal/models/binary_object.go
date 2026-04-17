package models

type BinaryObject struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Base64    string `json:"base64" gorm:"type:text;not null"`
	Filename  string `json:"filename" gorm:"not null"`
	Extension string `json:"extension" gorm:"not null"`
}
