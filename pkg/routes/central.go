package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func CentralRoutes(app *fiber.App) {
	// Health check
	app.Get("/", controllers.HealthCheck)

	app.Get("/supabase-user", func(c *fiber.Ctx) error {
		utils.LoginSupabaseUser()

		return c.SendString("Supabase user login initiated")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)
}
