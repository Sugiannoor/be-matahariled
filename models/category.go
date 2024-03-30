package models

import (
	"time"
)

type Category struct {
	CategoryId int64     `gorm:"primaryKey" json:"category_id"`
	Category   string    `gorm:"type:varchar(255);index" json:"category" validate:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
