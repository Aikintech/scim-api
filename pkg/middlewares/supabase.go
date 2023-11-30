package middlewares

import (
	"context"
	"os"

	"github.com/aikintech/scim/pkg/definitions"
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
)

func SupaAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		supabaseURL := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_KEY")
		supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)
		userToken := c.Get("X-USER-TOKEN", "")

		if userToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Middleware: Unauthorized",
			})
		}

		ctx := context.Background()
		user, err := supabaseClient.Auth.User(ctx, userToken)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Code:    fiber.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		injectSupabaseUserIntoContext(user, c)

		return c.Next()
	}
}

func injectSupabaseUserIntoContext(user *supabase.User, c *fiber.Ctx) {
	c.Locals("user", user)
}
