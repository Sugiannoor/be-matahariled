package models

import (
	"time"
)

type Product struct {
	ProductId   int64     `gorm:"primaryKey" json:"product_id" form:"product_id"`
	Title       string    `gorm:"type:varchar(255);index" json:"title" form:"title" validate:"required"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	FileId      int64     `json:"file_id"`
	File        File      `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASADE" json:"file"`
	CategoryId  int64     `json:"category_id" form:"category_id"`
	Category    Category  `json:"category"`
}

type ProductResponse struct {
	ProductId   int64     `json:"product_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FileId      int64     `json:"file_id"`
	CategoryId  int64     `json:"category_id"`
	PathFile    string    `json:"path_file"`
	Category    string    `json:"category"`
}

type ProductCreateRequest struct {
	Title       string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	File        File   `json:"file" form:"file" validate:"required"`
	CategoryId  int64  `json:"category_id" form:"category_id" validate:"required"`
}

type ProductEditRequest struct {
	ProductId   int64  `json:"product_id" form:"product_id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	File        File   `json:"file" form:"file"`
	CategoryId  int64  `json:"category_id" form:"category_id" validate:"required"`
}
