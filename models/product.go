package models

import (
	"time"
)

type Product struct {
	ProductId   int64     `gorm:"primaryKey" json:"product_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"type:varchar(100)" json:"category" validate:"required"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	FileId      int64     `json:"file_id"`           // Menyimpan ID File sebagai kunci asing
	File        File      `gorm:"foreignKey:FileId"` // Hubungan belongs to dengan File
}
