package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"

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
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	fmt.Println("Connecting to database...")

	config := &gorm.Config{}

	if os.Getenv("APP_ENV") == "local" {
		config.Logger = newLogger
	}

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), config)

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	if os.Getenv("APP_ENV") == "local" {
		DB = db.Debug()
	} else {
		DB = db
	}

	fmt.Println("Database connection established")
}

func MigrateDB() {
	// err := DB.AutoMigrate(
	// 	&models.User{}, &models.Podcast{},
	// 	&models.Post{}, &models.Playlist{},
	// 	&models.PrayerRequest{}, &models.Like{},
	// 	&models.Comment{}, &models.Event{},
	// 	&models.UserToken{}, &models.VerificationCode{},
	// 	&models.UserDevice{},
	// )

	// if err != nil {
	// 	return
	// }
}
