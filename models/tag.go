package models

import "time"

type Tag struct {
	TagId     int64     `gorm:"primaryKey" json:"tag_id"`
	Tag       string    `gorm:"type:varchar(255);index" json:"tag" validate:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TagRequest struct {
	Tag string `gorm:"type:varchar(255)" json:"tag" validate:"required"`
}
