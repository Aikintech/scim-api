package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func LoadRoutes(app *fiber.App) {
	fileController := controllers.NewFileController()

	// Health check
	app.Get("/", controllers.NewMiscController().HealthCheck)

	// Jobs
	app.Post("/j/backup", middlewares.CronJobsMiddleware(), controllers.NewMiscController().BackupDatabase)
	app.Post("/j/seed-podcast", middlewares.CronJobsMiddleware(), controllers.NewMiscController().SeedPodcasts)
	app.Post("/j/seed-roles-permissions", middlewares.CronJobsMiddleware(), controllers.NewMiscController().SeedRolesAndPermissions)

	// Upload files
	app.Post("/files", fileController.UploadFile)
	app.Get("/files/:fileKey", fileController.GetFileURL)
	app.Delete("/files/:fileKey", fileController.DeleteFile)

	MountClientRoutes(app)
	MountBackOfficeRoutes(app)
}
