package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aikintech/scim/pkg/models"
	"github.com/spf13/viper"
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
			Colorful:                  true,          // Disable color
		},
	)

	fmt.Println("Connecting to database...")

	dsn := viper.GetString("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	DB = db

	fmt.Println("Database connection established")
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Podcast{}, &models.Post{}, &models.Playlist{}, &models.PrayerRequest{})

	if err != nil {
		return
	}
}
