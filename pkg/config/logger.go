package config

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func InitializeLogger() {
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

}
