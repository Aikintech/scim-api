package routes

import (
	"net/http"

	"github.com/aikintech/scim/pkg/controllers"
	"github.com/labstack/echo/v4"
)

func LoadCentralRoutes(app *echo.Echo) {
	// Health check
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

	// Auth routes
	auth := app.Group("/auth")
	auth.POST("/sign-in", controllers.SignIn)

	// User routes
	app.GET("/users", controllers.GetUsers)
}
