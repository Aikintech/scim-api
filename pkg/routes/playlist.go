package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountPlaylistRoutes(app *fiber.App) {
	playlist := app.Group("/playlists")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")

	// Routes
	playlistController := controllers.NewPlaylistController()

	playlist.Get("/", jwtAuthWare, playlistController.GetPlaylists)
	playlist.Post("/", jwtAuthWare, playlistController.CreatePlaylist)
	playlist.Get("/:playlistId", jwtAuthWare, playlistController.GetPlaylist)
	playlist.Patch("/:playlistId", jwtAuthWare, playlistController.UpdatePlaylist)
	playlist.Delete("/:playlistId", jwtAuthWare, playlistController.DeletePlaylist)
	playlist.Post("/:playlistId/podcasts", jwtAuthWare, playlistController.AddPlaylistPodcasts)
	playlist.Patch("/:playlistId/podcasts", jwtAuthWare, playlistController.DeletePlaylistPodcasts)
}
