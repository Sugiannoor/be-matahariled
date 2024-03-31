package models

import "time"

type Contract struct {
	ContractId  int64     `gorm:"primaryKey" json:"contract_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Description string    `gorm:"type:text" json:"description"`
	StartDate   time.Time `gorm:"type:date" json:"start_date" validate:"required"`
	EndDate     time.Time `gorm:"type:date" json:"end_date" validate:"required"`
	Products    []Product `gorm:"foreignKey:ContractId" json:"products" validate:"required"`
}

type ContractResponse struct {
	ContractId   int64     `json:"contract_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	ProductNames []string  `json:"product_names"`
}

func GetProductNames(products []Product) []string {
	productNames := make([]string, len(products))
	for i, product := range products {
		productNames[i] = product.Title
	}
	return productNames
}
