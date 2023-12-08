package middlewares

import (
	"errors"
	"os"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	jwtWare "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func JWTMiddleware(accessType string) fiber.Handler {
	return jwtWare.New(jwtWare.Config{
		SigningKey:  jwtWare.SigningKey{Key: []byte(os.Getenv("APP_KEY"))},
		ContextKey:  constants.JWT_CONTEXT_KEY,
		TokenLookup: "header:X-USER-TOKEN",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			userJwt := c.Locals(constants.JWT_CONTEXT_KEY).(*jwt.Token)
			claims := userJwt.Claims.(jwt.MapClaims)

			// Refresh token
			if accessType == "refresh" {
				if tokenType := claims["tokenType"].(string); tokenType != "refresh" {
					return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
						Message: "Invalid token type provided",
					})
				}
			}

			// Get user
			user := new(models.User)
			if result := database.DB.Model(&models.User{}).Where("id = ?", claims["sub"].(string)).First(&user); result.Error != nil {
				if errors.Is(result.Error, gorm.ErrRecordNotFound) {
					return c.Status(fiber.StatusUnauthorized).JSON(definitions.MessageResponse{
						Message: "Invalid token provided",
					})
				} else {
					return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
						Message: result.Error.Error(),
					})
				}
			}

			c.Locals(constants.USER_CONTEXT_KEY, user)

			return c.Next()
		},
	})
}

func CronJobsMiddleware() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("CRON_USERNAME"): os.Getenv("CRON_PASSWORD"),
		},
	})
}
