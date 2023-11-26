package routes

import "github.com/gofiber/fiber/v2"

func LoadRoutes(app *fiber.App) {
	CentralRoutes(app)
	ClientRoutes(app)
}
