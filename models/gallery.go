package models

import "time"

type Gallery struct {
	GalleryId    int64      `gorm:"primaryKey" json:"gallery_id"`
	Path         string     `gorm:"type:varchar(255)" json:"path"`
	Gallery_name string     `gorm:"type:varchar(255)" json:"gallery"`
	Size         string     `gorm:"type:varchar(255)" json:"size"`
	Format       string     `gorm:"type:varchar(10)" json:"format"`
	CreatedAt    *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ProductId    int64      `json:"product_id"`
}
