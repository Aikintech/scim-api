package middlewares

import (
	"context"
	"os"

	"github.com/aikintech/scim/pkg/definitions"
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		supabaseURL := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_KEY")
		supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)
		userToken := c.Get("Authorization", "")

		if userToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "Unauthorized",
			})
		}

		ctx := context.Background()
		user, err := supabaseClient.Auth.User(ctx, c.Get("X-USER-TOKEN"))

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Code:    fiber.StatusUnauthorized,
				Message: err.Error(),
			})
		}

		c.Locals("user", user)

		return c.Next()
	}
}
