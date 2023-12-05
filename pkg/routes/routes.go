package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/gofiber/fiber/v2"
)

func LoadRoutes(app *fiber.App) {
	fileController := controllers.NewFileController()

	// Health check
	app.Get("/", controllers.HealthCheck)

	// Upload files
	app.Post("/files/upload", fileController.UploadFile)
	app.Get("/files/get-url", fileController.GetFileURL)

	ClientRoutes(app)
	BackOfficeRoutes(app)
}
