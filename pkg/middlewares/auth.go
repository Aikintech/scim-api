package middlewares

import (
	"fmt"
	"os"

	"github.com/aikintech/scim/pkg/definitions"
	jwtWare "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return jwtWare.New(jwtWare.Config{
		SigningKey:  jwtWare.SigningKey{Key: []byte(os.Getenv("APP_KEY"))},
		TokenLookup: "header:X-USER-TOKEN",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Code:    fiber.StatusUnauthorized,
				Message: err.Error(),
			})
		},
	})
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Get("X-USER-TOKEN", "")
		fmt.Println("userToken", userToken)

		return c.Next()
	}
}
