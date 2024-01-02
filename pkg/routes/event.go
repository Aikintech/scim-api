package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountEventRoutes(app *fiber.App) {
	events := app.Group("/events")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")

	events.Get("/my-calendar-events", jwtAuthWare, controllers.NewEventController().MyCalendarEvents)
	events.Get("/", controllers.NewEventController().GetEvents)
	events.Get("/:eventId", controllers.NewEventController().GetEvent)
	events.Patch("/:eventId", jwtAuthWare, controllers.NewEventController().SyncEventToCalendar)
}
