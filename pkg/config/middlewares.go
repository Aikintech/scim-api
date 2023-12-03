package config

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
)

func LoadGlobalMiddlewares(app *fiber.App) {
	app.Use(cors.New())
	app.Use(limiter.New())

	if os.Getenv("APP_ENV") == "local" {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger: &logger,
		}))
	}
}
