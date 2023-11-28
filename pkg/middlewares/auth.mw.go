package middlewares

import (
	"context"

	"github.com/aikintech/scim/pkg/definitions"
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/spf13/viper"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		supabaseURL := viper.GetString("SUPABASE_URL")
		supabaseKey := viper.GetString("SUPABASE_KEY")
		supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)

		ctx := context.Background()
		user, err := supabaseClient.Auth.User(ctx, c.Get("Authorization"))

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
