package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func MountTestimonyRoutes(app *fiber.App) {
	testimony := app.Group("/testimonies")

	testimonyController := controllers.NewTestimonyController()

	// Routes
	testimony.Get("/", testimonyController.GetTestimonies)
	testimony.Get("/:testimonyId", testimonyController.GetTestimony)
}
