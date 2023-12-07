package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	app.Get("/home", controllers.NewHomeController().ClientHome)

	// Routers
	MountAuthRoutes(app)
	MountPodcastRoutes(app)
	MountPodcastRoutes(app)
	MountPrayerRoutes(app)
}
