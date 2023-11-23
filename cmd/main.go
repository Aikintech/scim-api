package main

import (
	"os"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/routes"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
)

func init() {
	// Load environment variables
	// env.LoadEnv()

	// Load database
	config.ConnectDB()
	config.MigrateDB()

	// Configure validation
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = false
	})
}

func main() {
	// Instantiate the app
	app := fiber.New()

	// Global middleware
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// Load routes
	routes.LoadRoutes(app)

	// Start the app
	if err := app.Listen(":9000"); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
