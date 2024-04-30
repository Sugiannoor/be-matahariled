package initialize

import (
	"Matahariled/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/matahariled?parseTime=true"))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Contract{})
	db.AutoMigrate(&models.History{})
	db.AutoMigrate(&models.Video{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.Gallery{})
	DB = db
}
