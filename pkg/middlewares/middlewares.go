package middlewares

import (
	"os"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/limiter"
)

func LoadGlobalMiddlewares(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
	}))

	// app.Use(limiter.New(limiter.Config{
	// 	Max: 60,
	// }))

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &config.Logger,
	}))
}
