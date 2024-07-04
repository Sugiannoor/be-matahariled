package initialize

import (
	"Matahariled/config"
	"Matahariled/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
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
