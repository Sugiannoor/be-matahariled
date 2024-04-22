package models

import "time"

type Video struct {
	VideoId   int64     `gorm:"primaryKey" json:"video_id"`
	Title     string    `gorm:"type:varchar(255);index" json:"video_title" validate:"required"`
	Embed     string    `gorm:"type:varchar(255);index" json:"embed" validate:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
