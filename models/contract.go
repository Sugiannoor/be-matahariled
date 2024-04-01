package models

type Contract struct {
	ContractId  int64     `gorm:"primaryKey" json:"contract_id"`
	Title       string    `gorm:"type:varchar(255)" json:"title" validate:"required"`
	Description string    `gorm:"type:text" json:"description"`
	StartDate   string    `gorm:"type:varchar(20)" json:"start_date" validate:"required"`
	EndDate     string    `gorm:"type:varchar(20)" json:"end_date" validate:"required"`
	UserID      int64     `json:"user_id" form:"user_id" validate:"required"`
	User        User      `json:"user"`
	Products    []Product `gorm:"many2many:ContractProduct;constraint:OnDelete:CASCADE" json:"products"`
}

type ContractProduct struct {
	ContractId int64 `gorm:"primaryKey"`
	ProductId  int64 `gorm:"primaryKey"`
}

type ContractResponse struct {
	ContractId   int64    `json:"contract_id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	StartDate    string   `gorm:"type:varchar(20)" json:"start_date"`
	EndDate      string   `gorm:"type:varchar(20)" json:"end_date"`
	UserName     string   `json:"username"`
	ProductNames []string `json:"product"`
}

type ContractRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	StartDate   string  `gorm:"type:varchar(20)" json:"start_date" validate:"required"`
	EndDate     string  `gorm:"type:varchar(20)" json:"end_date" validate:"required"`
	UserID      int64   `json:"user_id" validate:"required"`
	ProductIDs  []int64 `json:"product_ids" validate:"required"`
}

func GetProductNames(products []Product) []string {
	productNames := make([]string, len(products))
	for i, product := range products {
		productNames[i] = product.Title
	}
	return productNames
}
