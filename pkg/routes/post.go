package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountPostRoutes(app *fiber.App) {
	posts := app.Group("/posts")
	jwtAuthWare := middlewares.JWTMiddleware("access")

	postController := controllers.NewPostController()

	posts.Get("/", jwtAuthWare, postController.GetPosts)
	posts.Get("/:postId", jwtAuthWare, postController.GetPost)
	posts.Get("/:postId/comments", postController.GetPostComments)
	posts.Post("/:postId/comments", jwtAuthWare, postController.CreatePostComment)
	posts.Patch("/:postId/comments/:commentId", jwtAuthWare, postController.UpdatePostComment)
	posts.Delete("/:postId/comments/:commentId", jwtAuthWare, postController.DeletePostComment)
}
