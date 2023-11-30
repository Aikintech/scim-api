package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/aikintech/scim/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CentralRoutes(app *fiber.App) {
	// Health check
	app.Get("/", controllers.HealthCheck)

	// Middlewares
	refreshJwtAuthWare := middlewares.JWTMiddleware("refresh")

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)
	auth.Get("/refresh-token", refreshJwtAuthWare, controllers.RefreshToken)
}
