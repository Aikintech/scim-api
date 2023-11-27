package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	// Create a new sub-router (client)
	client := app.Group("/client")

	// Podcasts
	client.Get("/Podcasts", controllers.ClientGetPodcasts)
}
