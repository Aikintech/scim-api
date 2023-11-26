package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func CentralRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)
}
