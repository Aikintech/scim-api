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

	// TODO: Add comments and likes to the cache
	// listAllPodcastsCache := middlewares.AllPodcastsCache()

	// TODO: Add comments and likes to the cache
	// podcastByIdCache := middlewares.PodcastByIdCache()

	// Routes
	podcastController := controllers.NewPodcastController()
	commentController := controllers.NewCommentController()
	likeController := controllers.NewLikeController()

	podcasts.Get("/", podcastController.ListPodcasts)
	podcasts.Get("/all", podcastController.ListPodcasts)
	podcasts.Get("/:podcastId", podcastController.ShowPodcast)
	podcasts.Get("/:podcastId/comments", commentController.GetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuthWare, commentController.StorePodcastComment)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuthWare, commentController.UpdatePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuthWare, likeController.LikePodcast)
	podcasts.Delete("/:podcastId/comments/:commentId", jwtAuthWare, commentController.DeletePodcastComment)
}
