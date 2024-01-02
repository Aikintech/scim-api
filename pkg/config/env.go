package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		if os.IsNotExist(err) {
			Logger.Error().Msg(".env file not found, using only system environment variables.")
		} else {
			Logger.Error().Msgf("error loading .env file: %v", err)
		}
	}
}
