package initialize

import (
	"Matahariled/models"
	"fmt"
	"os"

	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDSN() string {
	gotenv.Load()
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, dbName)
}

func ConnectDatabase() {
	dsn := GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
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
	db.AutoMigrate(&models.Hero{})
	DB = db
}
