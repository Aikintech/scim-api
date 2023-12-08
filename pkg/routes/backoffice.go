package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BackOfficeRoutes(app *fiber.App) {
	backoffice := app.Group("/backoffice")
	events := backoffice.Group("/events")
	prayers := backoffice.Group("/prayer-requests")
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Controller initializations
	prayerController := controllers.NewPrayerController()
	eventController := controllers.NewEventController()

	// Events
	events.Post("/", jwtAuthWare, eventController.BackofficeCreateEvent)

	// Prayer requests
	prayers.Get("/", prayerController.BackOfficeGetPrayers)
}
