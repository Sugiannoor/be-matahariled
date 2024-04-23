package models

import (
	"time"
)

type Category struct {
	CategoryId int64     `gorm:"primaryKey" json:"category_id"`
	Category   string    `gorm:"type:varchar(255);index" json:"category" validate:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Tags       []Tag     `gorm:"many2many:CategoryTag;constraint:OnDelete:CASCADE" json:"tags"`
}

type CategoryTag struct {
	CategoryId int64 `gorm:"primaryKey"`
	TagId      int64 `gorm:"primaryKey"`
}

type CategoryRequest struct {
	Category string  `gorm:"type:varchar(255);index" json:"category" validate:"required"`
	TagIDs   []int64 `json:"tag_ids"`
}
