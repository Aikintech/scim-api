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

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware()

	// Podcasts
	podcasts.Post("/seed", controllers.SeedPodcasts)

	podcasts.Get("/", controllers.ClientListPodcasts)
	podcasts.Get("/:podcastId", controllers.ClientShowPodcast)
	podcasts.Get("/:podcastId/comments", controllers.ClientGetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuthWare, controllers.ClientStorePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuthWare, controllers.ClientLikePodcast)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuthWare, controllers.ClientUpdatePodcastComment)

	// Playlists
	playlists.Post("/", jwtAuthWare, controllers.ClientCreatePlaylist)

	// Prayer requests
	prayerRequests.Post("/", controllers.ClientRequestPrayer)
}
