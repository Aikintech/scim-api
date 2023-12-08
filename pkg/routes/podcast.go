package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountPodcastRoutes(app *fiber.App) {
	podcasts := app.Group("/podcasts")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")
	// podcastsCache := middlewares.PodcastsCache()
	listAllPodcastsCache := middlewares.AllPodcastsCache()
	podcastByIdCache := middlewares.PodcastByIdCache()

	// Routes
	podcastController := controllers.NewPodcastController()
	commentController := controllers.NewCommentController()
	likeController := controllers.NewLikeController()

	podcasts.Get("/", podcastController.ListPodcasts)
	podcasts.Get("/all", listAllPodcastsCache, podcastController.ListAllPodcasts)
	podcasts.Get("/:podcastId", podcastByIdCache, podcastController.ShowPodcast)
	podcasts.Get("/:podcastId/comments", commentController.GetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuthWare, commentController.StorePodcastComment)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuthWare, commentController.UpdatePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuthWare, likeController.LikePodcast)
	podcasts.Delete("/:podcastId/comments/:commentId", jwtAuthWare, commentController.DeletePodcastComment)
}
