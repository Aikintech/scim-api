package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func LoadRoutes(app *fiber.App) {
	// Health check
	app.Get("/", controllers.HealthCheck)

	ClientRoutes(app)
}
