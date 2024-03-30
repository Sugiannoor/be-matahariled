package models

import "time"

type User struct {
	UserId      int64      `gorm:"primaryKey" json:"user_id"`
	FullName    string     `gorm:"type:varchar(100);index" json:"full_name" validate:"required"`
	UserName    string     `gorm:"type:varchar(100);index" json:"username" validate:"required" `
	PhoneNumber string     `gorm:"type:varchar(100)" json:"phone_number" validate:"required" `
	Password    string     `gorm:"type:varchar(100)" json:"password" validate:"required"`
	Email       string     `gorm:"type:varchar(100)" json:"email" validate:"required"`
	Address     *string    `gorm:"type:varchar(300)" json:"address"`
	Role        string     `gorm:"type:ENUM('Admin', 'Customer', 'SuperAdmin'); default:'Customer'" json:"role"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
