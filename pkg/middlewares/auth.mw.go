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

		injectUserIntoContext(user, c)

		return c.Next()
	}
}

func injectUserIntoContext(user *supabase.User, c *fiber.Ctx) {
	c.Locals("user", user)
}
