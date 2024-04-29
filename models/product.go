package models

import (
	"time"
)

type Product struct {
	ProductId     int64     `gorm:"primaryKey" json:"product_id" form:"product_id"`
	Title         string    `gorm:"type:varchar(255);index" json:"title" form:"title" validate:"required"`
	Description   string    `gorm:"type:text" json:"description"`
	Specification string    `gorm:"type:text" json:"specification"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	FileId        int64     `json:"file_id"`
	File          File      `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASCADE" json:"file"`
	CategoryId    int64     `json:"category_id" form:"category_id"`
	Category      Category  `json:"category"`
	Galleries     []Gallery `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASCADE" json:"galleries"`
}

type ProductResponse struct {
	ProductId     int64     `json:"product_id"`
	Title         string    `json:"title"`
	Specification string    `json:"specification"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	FileId        int64     `json:"file_id"`
	CategoryId    int64     `json:"category_id"`
	PathFile      string    `json:"path_file"`
	Category      string    `json:"category"`
}
