package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken := c.Get("X-USER-TOKEN", "")
		fmt.Println("userToken", userToken)

		return c.Next()
	}
}
