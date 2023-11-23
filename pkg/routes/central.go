package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func LoadCentralRoutes(app *fiber.App) {
	// Health check
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/sign-in", controllers.SignIn)

	// User routes
	app.Get("/users", controllers.GetUsers)
}
