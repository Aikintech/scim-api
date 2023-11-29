package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		// It's okay if the .env file is not present
		if os.IsNotExist(err) {
			fmt.Println(".env file not found, using only system environment variables.")
		} else {
			fmt.Errorf("error loading .env file: %v", err)
		}
	}
}
