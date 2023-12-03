package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aikintech/scim/pkg/routes"

	"github.com/aikintech/scim/pkg/config"
	"github.com/gofiber/fiber/v2"
)

func init() {
	// Load environment variables
	config.LoadEnv()

	// Initialize redis
	config.InitializeRedis()

	// Load database
	config.ConnectDB()
	config.MigrateDB()
}

func main() {
	// Instantiate a new fiber app
	app := fiber.New()

	// Middlewares
	config.LoadGlobalMiddlewares(app)

	// Routes
	routes.LoadRoutes(app)

	// Start the app
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatal(err.Error())
	}
}
