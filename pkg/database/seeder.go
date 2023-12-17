package database

import (
	"fmt"
	"os"
)

func RunDatabaseSeeder(run bool) {
	if os.Getenv("APP_ENV") == "local" {
		// Run seeder locally
		fmt.Println(run)
	}
}
