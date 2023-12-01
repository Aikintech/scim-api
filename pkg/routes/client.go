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

	/**
	*** Auth routes
	**/
	auth := app.Group("/auth")
	auth.Post("/login", controllers.Login)
	auth.Post("/register", controllers.Register)
	auth.Post("/forgot-password", controllers.ForgotPassword)
	auth.Post("/resend-email-verification", controllers.ResendEmailVerification)
	auth.Get("/refresh-token", refreshJwtAuthWare, controllers.RefreshToken)

	/**
	*** Podcast and playlist routes
	**/
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

	/**
	*** Prayer request
	**/
	prayerRequests.Post("/", controllers.ClientRequestPrayer)
}
