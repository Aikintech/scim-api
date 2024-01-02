package main

import (
	"fmt"
	"os"

	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/middlewares"

	"github.com/aikintech/scim-api/pkg/routes"
	"github.com/aikintech/scim-api/pkg/utils"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/gofiber/fiber/v2"
)

func init() {
	// Load environment variables
	config.LoadEnv()

	// Initialize logger
	config.InitializeLogger()

	// Initialize redis
	config.InitializeRedis()

	// Load database
	database.ConnectDB()
	database.MigrateDB()
}

func main() {
	// Instantiate a new fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 64 * 1024 * 1024, // 64MB
	})

	// Middlewares
	middlewares.LoadGlobalMiddlewares(app)

	// Routes
	routes.LoadRoutes(app)

	// Dump routes to a file
	if os.Getenv("APP_ENV") == "local" {
		if err := utils.DumpRoutesToFile(app); err != nil {
			config.Logger.Error().Msgf("Loading .env variables: %s", err.Error())
		}
	}

	// Start the app
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		config.Logger.Fatal().Msgf("Starting fiber application: %s", err.Error())
	}
}
