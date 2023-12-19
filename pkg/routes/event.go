package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func MountEventRoutes(app *fiber.App) {
	events := app.Group("/events")

	events.Get("/", controllers.NewEventController().GetEvents)
	events.Get("/:eventId", controllers.NewEventController().GetEvent)
}
