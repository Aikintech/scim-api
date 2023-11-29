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

	// Podcasts
	podcasts.Post("/seed", controllers.SeedPodcasts)

	podcasts.Get("/", controllers.ClientListPodcasts)
	podcasts.Get("/:podcastId", controllers.ClientShowPodcast)
	podcasts.Patch("/:podcastId/like", middlewares.AuthMiddleware(), controllers.ClientLikePodcast)
	podcasts.Post("/:podcastId/comment", middlewares.AuthMiddleware(), controllers.ClientCommentPodcast)
	podcasts.Patch("/:podcastId/comment/:commentId", middlewares.AuthMiddleware(), controllers.ClientUpdatePodcastComment)

	// Playlists
	playlists.Post("/", controllers.ClientCreatePlaylist)

	// Prayer requests
	prayerRequests.Post("/", controllers.ClientRequestPrayer)
}
