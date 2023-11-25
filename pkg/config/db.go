package config

import (
	"log"

	"github.com/aikintech/scim/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// dsn := "avnadmin:AVNS_H7JBDGn1LLrlQXCB84T@tcp(scim-app-test-dabeixin-8a11.a.aivencloud.com:14958)/scim-test?charset=utf8mb4&parseTime=True"
	dsn := "dbuser:Password2020!@tcp(127.0.0.1:3306)/chms?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	DB = db
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Podcast{})

	if err != nil {
		return
	}
}
