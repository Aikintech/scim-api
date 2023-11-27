package config

import (
	"log"
	"os"
	"time"

	"github.com/aikintech/scim/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	// dsn := "avnadmin:AVNS_H7JBDGn1LLrlQXCB84T@tcp(scim-app-test-dabeixin-8a11.a.aivencloud.com:14958)/scim-test?charset=utf8mb4&parseTime=True"
	dsn := "dbuser:Password2020!@tcp(127.0.0.1:3306)/chms?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

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
