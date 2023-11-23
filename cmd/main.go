package main

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/routes"
	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	app := echo.New()

	// Global middleware
	app.Use(middleware.Logger())

	// Load routes
	routes.LoadRoutes(app)

	// Start the app
	// app.Logger.Fatal(app.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
	app.Logger.Fatal(app.Start(":9000"))

	println("Hello, World!")
}
