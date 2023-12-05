package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ClientRoutes(app *fiber.App) {
	// Groups
	auth := app.Group("/auth")
	podcasts := app.Group("/podcasts")
	playlists := app.Group("/playlists")
	prayers := app.Group("/prayer-requests")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")
	refreshJwtAuthWare := middlewares.JWTMiddleware("refresh")
	// podcastsCache := middlewares.PodcastsCache()
	listAllPodcastsCache := middlewares.AllPodcastsCache()
	podcastByIdCache := middlewares.PodcastByIdCache()

	/**
	*** Auth routes
	**/
	authController := controllers.NewAuthController()

	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/resend-email-verification", authController.ResendEmailVerification)
	auth.Get("/refresh-token", refreshJwtAuthWare, authController.RefreshToken)

	/**
	*** Podcast and playlist routes
	**/
	// Podcasts
	podcastController := controllers.NewPodcastController()
	commentController := controllers.NewCommentController()
	likeController := controllers.NewLikeController()

	podcasts.Post("/seed", middlewares.CronJobsMiddleware(), podcastController.SeedPodcasts)
	podcasts.Get("/", podcastController.ListPodcasts)
	podcasts.Get("/all", listAllPodcastsCache, podcastController.ListAllPodcasts)
	podcasts.Get("/:podcastId", podcastByIdCache, podcastController.ShowPodcast)
	podcasts.Get("/:podcastId/comments", commentController.GetPodcastComments)
	podcasts.Post("/:podcastId/comments", jwtAuthWare, commentController.StorePodcastComment)
	podcasts.Patch("/:podcastId/comments/:commentId", jwtAuthWare, commentController.UpdatePodcastComment)
	podcasts.Patch("/:podcastId/like", jwtAuthWare, likeController.LikePodcast)
	podcasts.Delete("/:podcastId/comments/:commentId", jwtAuthWare, commentController.DeletePodcastComment)

	// Playlists
	playlistController := controllers.NewPlaylistController()

	playlists.Get("/", jwtAuthWare, playlistController.GetPlaylists)
	playlists.Post("/", jwtAuthWare, playlistController.CreatePlaylist)
	playlists.Get("/:playlistId", jwtAuthWare, playlistController.GetPlaylist)
	playlists.Patch("/:playlistId", jwtAuthWare, playlistController.UpdatePlaylist)
	playlists.Delete("/:playlistId", jwtAuthWare, playlistController.DeletePlaylist)

	/**
	*** Prayer request
	**/
	prayerController := controllers.NewPrayerController()

	prayers.Get("/", jwtAuthWare, prayerController.MyPrayers)
	prayers.Post("/", jwtAuthWare, prayerController.RequestPrayer)

	// Dashboard/home

	app.Get("/home", controllers.NewHomeController().ClientHome)
}
