package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func BackOfficeRoutes(app *fiber.App) {
	backoffice := app.Group("/backoffice")
	events := backoffice.Group("/events")
	posts := backoffice.Group("/posts")
	prayers := backoffice.Group("/prayer-requests")
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Controller initializations
	prayerController := controllers.NewPrayerController()
	eventController := controllers.NewEventController()
	postController := controllers.NewPostController()

	// Events
	events.Get("/", jwtAuthWare, eventController.BackofficeGetEvents)
	events.Post("/", jwtAuthWare, eventController.BackofficeCreateEvent)
	events.Get("/:eventId", jwtAuthWare, eventController.BackofficeGetEvent)
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
}
