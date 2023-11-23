package routes

import (
	"github.com/gofiber/fiber/v2"
)

func LoadRoutes(app *fiber.App) {
	// Central routes
	LoadCentralRoutes(app)

	// Client routes
	LoadClientRoutes(app)
}
