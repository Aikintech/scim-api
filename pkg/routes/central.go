package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func CentralRoutes(app *fiber.App) {
	// Health check
	app.Get("/", controllers.HealthCheck)
}
