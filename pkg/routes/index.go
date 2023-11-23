package routes

import "github.com/labstack/echo/v4"

func LoadRoutes(app *echo.Echo) {
	// Central routes
	LoadCentralRoutes(app)

	// Client routes
	LoadClientRoutes(app)
}
