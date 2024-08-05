package models

import (
	"time"
)

type User struct {
	UserId      int64      `gorm:"primaryKey" json:"user_id"`
	FullName    string     `gorm:"type:varchar(100);index" json:"full_name"`
	UserName    string     `gorm:"type:varchar(100);index" json:"username" `
	PhoneNumber string     `gorm:"type:varchar(100)" json:"phone_number" `
	Password    string     `gorm:"type:varchar(100)" json:"password"`
	Email       string     `gorm:"type:varchar(100)" json:"email"`
	Address     *string    `gorm:"type:varchar(300)" json:"address"`
	Role        string     `gorm:"type:ENUM('Admin', 'Customer', 'SuperAdmin'); default:'Customer'" json:"role"`
	FileId      int64      `json:"file_id"`
	File        File       `gorm:"constraint:OnDelete:CASCADE;OnUpdate:CASADE" json:"file"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
type UserResponse struct {
	UserId      int64      `gorm:"primaryKey" json:"user_id"`
	FullName    string     `gorm:"type:varchar(100);index" json:"full_name" validate:"required"`
	UserName    string     `gorm:"type:varchar(100);index" json:"username" validate:"required" `
	PhoneNumber string     `gorm:"type:varchar(100)" json:"phone_number" validate:"required" `
	Password    string     `gorm:"-"`
	Email       string     `gorm:"type:varchar(100)" json:"email" validate:"required"`
	Address     *string    `gorm:"type:varchar(300)" json:"address"`
	Role        string     `gorm:"type:ENUM('Admin', 'Customer', 'Superadmin'); default:'Customer'" json:"role"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
