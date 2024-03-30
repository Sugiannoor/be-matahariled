package initialize

import (
	"Matahariled/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/db_matahariled?parseTime=true"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&models.Product{})
	DB = db
}
