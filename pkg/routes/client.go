package routes

import (
	"github.com/aikintech/scim/pkg/controllers"
	"github.com/aikintech/scim/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	// Groups
	podcasts := app.Group("/podcasts")
	playlists := app.Group("/playlists")
	prayerRequests := app.Group("/prayer-requests")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")
	refreshJwtAuthWare := middlewares.JWTMiddleware("refresh")
	// podcastsCache := middlewares.PodcastsCache()
	listAllPodcastsCache := middlewares.AllPodcastsCache()
	podcastByIdCache := middlewares.PodcastByIdCache()

	/**
	*** Auth routes
	**/
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/reset-password", controllers.ResetPassword)
	auth.Post("/resend-email-verification", controllers.ResendEmailVerification)
	auth.Get("/refresh-token", refreshJwtAuthWare, controllers.RefreshToken)

	/**
	*** Podcast and playlist routes
	**/
	// Podcasts
	podcasts.Post("/seed", controllers.SeedPodcasts)
	podcasts.Get("/", controllers.ListPodcasts)
	podcasts.Get("/all", listAllPodcastsCache, controllers.ListAllPodcasts)
	podcasts.Get("/:podcastId", podcastByIdCache, controllers.ShowPodcast)
	podcasts.Get("/:podcastId/comments", controllers.GetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuthWare, controllers.StorePodcastComment)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuthWare, controllers.UpdatePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuthWare, controllers.LikePodcast)
	podcasts.Delete("/:podcastId/comments/:commentId", jwtAuthWare, controllers.DeletePodcastComment)

	// Playlists
	playlists.Get("/", jwtAuthWare, controllers.GetPlaylists)
	playlists.Post("/", jwtAuthWare, controllers.CreatePlaylist)
	playlists.Get("/:playlistId", jwtAuthWare, controllers.GetPlaylist)
	playlists.Patch("/:playlistId", jwtAuthWare, controllers.UpdatePlaylist)
	playlists.Delete("/:playlistId", jwtAuthWare, controllers.DeletePlaylist)

	// Comments

	/**
	*** Prayer request
	**/
	prayerRequests.Post("/", controllers.RequestPrayer)
}
