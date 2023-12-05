package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aikintech/scim/pkg/routes"
	"github.com/aikintech/scim/pkg/utils"

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
}

func main() {
	// Instantiate a new fiber app
	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: 64 * 1024 * 1024, // 64MB
	})

	if !fiber.IsChild() {
		config.MigrateDB()
	}

	// Middlewares
	config.LoadGlobalMiddlewares(app)

	// Routes
	routes.LoadRoutes(app)

	// Dump routes to a file
	if os.Getenv("APP_ENV") == "local" {
		if err := utils.DumpRoutesToFile(app); err != nil {
			fmt.Println("Error:", err.Error())
		}
	}

	// Start the app
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatal(err.Error())
	}
}
