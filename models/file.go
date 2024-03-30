package models

import "time"

type File struct {
	FileId    int64      `gorm:"primaryKey" json:"file_id"`
	Path      string     `gorm:"type:varchar(255)" json:"path"`
	File_name string     `gorm:"type:varchar(255)" json:"file_name"`
	Size      string     `gorm:"type:varchar(255)" json:"size"`
	Format    string     `gorm:"type:varchar(10)" json:"format"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
