package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func MountClientRoutes(app *fiber.App) {
	app.Get("/home", controllers.NewHomeController().ClientHome)

	// Routers
	MountAuthRoutes(app)
	MountPodcastRoutes(app)
	MountPlaylistRoutes(app)
	MountPrayerRoutes(app)
	MountPostRoutes(app)
	MountTransactionRoutes(app)
	MountEventRoutes(app)
}
