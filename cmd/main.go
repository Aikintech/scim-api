package main

import (
	"fmt"
	"github.com/aikintech/scim/pkg/routes"
	"os"

	"github.com/aikintech/scim/pkg/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func init() {
	// Load environment variables
	// env.LoadEnv()

	// Load database
	config.ConnectDB()
	config.MigrateDB()
}

func main() {
	// Instantiate the app
	app := fiber.New()

	// Global middleware
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
