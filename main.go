package main

import (
	"fmt"
	"os"

	"github.com/aikintech/scim/pkg/routes"

	"github.com/aikintech/scim/pkg/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func init() {
	// Load environment variables
	config.LoadEnv()

	// Load database
	config.ConnectDB()
	config.MigrateDB()
}

func main() {
	// Instantiate a new fiber app
	app := fiber.New()

	// Global middlewares
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// Routes
	routes.LoadRoutes(app)

	// Start the app
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
