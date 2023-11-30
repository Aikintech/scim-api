package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/aikintech/scim/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	// Create a new sub-router (client)
	client := app.Group("/client")
	podcasts := client.Group("/podcasts")
	playlists := client.Group("/playlists")
	prayerRequests := client.Group("/prayer-requests")
	jwtAuth := middlewares.AuthMiddleware()

	// Podcasts
	podcasts.Post("/seed", controllers.SeedPodcasts)

	podcasts.Get("/", controllers.ClientListPodcasts)
	podcasts.Get("/:podcastId", controllers.ClientShowPodcast)
	podcasts.Get("/:podcastId/comments", controllers.ClientGetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuth, controllers.ClientStorePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuth, controllers.ClientLikePodcast)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuth, controllers.ClientUpdatePodcastComment)

	// Playlists
	playlists.Post("/", jwtAuth, controllers.ClientCreatePlaylist)

	// Prayer requests
	prayerRequests.Post("/", controllers.ClientRequestPrayer)
}
