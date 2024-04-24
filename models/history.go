package models

import "time"

type History struct {
	HistoryId   int64     `gorm:"primaryKey" json:"history_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	StartDate   string    `gorm:"type:varchar(20)" json:"start_date"`
	EndDate     string    `gorm:"type:varchar(20)" json:"end_date"`
	ProductId   int64     `gorm:"index" json:"product_id"`
	Product     Product   `json:"product"`
	FileId      int64     `gorm:"index" json:"file_id"`
	File        File      `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASCADE" json:"file"`
	VideoId     int64     `gorm:"index" json:"video_id"`
	Video       Video     `gorm:"constraint:onDelete:CASCADE;OnUpdate:CASCADE" json:"video"`
	UserId      int64     `gorm:"index" json:"user_id"`
	User        User      `gorm:"constraint:onDelete:CASCADE;OnUpdate:CASCADE" json:"user"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type HistoryResponse struct {
	HistoryId    int64     `json:"history_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartDate    string    `gorm:"type:varchar(20)" json:"start_date"`
	EndDate      string    `gorm:"type:varchar(20)" json:"end_date"`
	ProductName  string    `json:"product"`
	CategoryName string    `json:"category"`
	PathFile     string    `json:"path_file"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
