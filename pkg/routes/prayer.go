package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountPrayerRoutes(app *fiber.App) {
	prayers := app.Group("/prayers")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Routes
	prayerController := controllers.NewPrayerController()

	prayers.Get("/", jwtAuthWare, prayerController.MyPrayers)
	prayers.Post("/", jwtAuthWare, prayerController.RequestPrayer)
}
