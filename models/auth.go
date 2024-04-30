package models

type LoginRequest struct {
	Email    string `gorm:"type:varchar(100)" json:"email" validate:"required"`
	Password string `gorm:"type:varchar(100)" json:"password" validate:"required"`
}
