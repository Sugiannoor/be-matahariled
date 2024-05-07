package models

import "time"

type Hero struct {
	HeroId    int64      `gorm:"primaryKey" json:"hero_id"`
	Path      string     `gorm:"type:varchar(255)" json:"path"`
	Hero_name string     `gorm:"type:varchar(255)" json:"hero_name"`
	Size      string     `gorm:"type:varchar(255)" json:"size"`
	Format    string     `gorm:"type:varchar(10)" json:"format"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ProductId   int64     `gorm:"index" json:"product_id"`
	Product     Product   `json:"product"`
}