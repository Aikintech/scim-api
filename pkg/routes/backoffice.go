package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BackOfficeRoutes(app *fiber.App) {
	backoffice := app.Group("/backoffice")
	events := backoffice.Group("/events")
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Events
	eventController := controllers.NewEventController()

	events.Post("/", jwtAuthWare, eventController.BackofficeStoreEvent)
}
