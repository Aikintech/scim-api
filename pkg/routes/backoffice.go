package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountBackOfficeRoutes(app *fiber.App) {
	backoffice := app.Group("/backoffice")
	events := backoffice.Group("/events")
	posts := backoffice.Group("/posts")
	prayers := backoffice.Group("/prayer-requests")
	transactions := backoffice.Group("/transactions")
	users := backoffice.Group("/users")
	testimonies := backoffice.Group("/testimonies")

	// TODO change middleware to rbacAuthMiddleware
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Controller initializations
	prayerController := controllers.NewPrayerController()
	eventController := controllers.NewEventController()
	postController := controllers.NewPostController()
	transactionController := controllers.NewTransactionController()
	userController := controllers.NewUserController()
	fileController := controllers.NewFileController()
	testimonyController := controllers.NewTestimonyController()

	// Events
	events.Get("/", jwtAuthWare, eventController.BackofficeGetEvents)
	events.Post("/", jwtAuthWare, eventController.BackofficeCreateEvent)
	events.Get("/:eventId", jwtAuthWare, eventController.GetEvent)
	events.Patch("/:eventId", jwtAuthWare, eventController.BackofficeUpdateEvent)
	events.Delete("/:eventId", jwtAuthWare, eventController.BackofficeDeleteEvent)

	// Prayer requests
	prayers.Get("/", jwtAuthWare, prayerController.BackOfficeGetPrayers)

	// Posts
	posts.Get("/", jwtAuthWare, postController.GetPosts)
	posts.Post("/", jwtAuthWare, postController.BackofficeCreatePost)
	posts.Get("/:postId", jwtAuthWare, postController.GetPost)
	posts.Patch("/:postId", jwtAuthWare, postController.BackofficeUpdatePost)
	posts.Delete("/:postId", jwtAuthWare, postController.BackofficeDeletePost)

	// Transactions
	transactions.Get("/", jwtAuthWare, transactionController.BackofficeGetTransactions)
	transactions.Get("/:transactionId", jwtAuthWare, transactionController.BackofficeGetTransaction)

	// Users
	users.Get("/", jwtAuthWare, userController.BackofficeGetUsers)

	// Testimonies
	testimonies.Get("/", jwtAuthWare, testimonyController.BackofficeGetTestimonies)

	// Files
	backoffice.Post("/files", fileController.UploadFile)
	backoffice.Delete("/files", jwtAuthWare, fileController.DeleteFile)
}
