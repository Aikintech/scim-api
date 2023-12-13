package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")
	refreshJwtAuthWare := middlewares.JWTMiddleware("refresh")

	// Routes
	authController := controllers.NewAuthController()

	auth.Post("/login", authController.Login)
	auth.Post("/register", authController.Register)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/resend-email-verification", authController.ResendEmailVerification)
	auth.Post("/verify-account", authController.VerifyAccount)

	auth.Get("/refresh-token", refreshJwtAuthWare, authController.RefreshToken)
	auth.Get("/user", jwtAuthWare, authController.User)

}
