package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func MountAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")

	// Middlewares
	jwtAuthWare := middlewares.JWTMiddleware("access")
	refreshJwtAuthWare := middlewares.JWTMiddleware("refresh")

	// Routes
	authController := controllers.NewAuthController()

	// Public routes
	auth.Post("/login", limiter.New(limiter.Config{Max: 60}), authController.Login)
	auth.Post("/register", authController.Register)
	auth.Post("/forgot-password", authController.ForgotPassword)
	auth.Patch("/reset-password", authController.ResetPassword)
	auth.Post("/resend-email-verification", authController.ResendEmailVerification)
	auth.Patch("/verify-account", authController.VerifyAccount)
	auth.Post("/verify-code", authController.VerifyCode)
	auth.Post("/social", authController.SocialAuth)

	// Protected routes
	auth.Get("/refresh-token", refreshJwtAuthWare, authController.RefreshToken)
	auth.Get("/user", jwtAuthWare, authController.User)

	auth.Patch("/update-user-avatar", jwtAuthWare, authController.UpdateUserAvatar)
	auth.Patch("/update-user-details", jwtAuthWare, authController.UpdateUserDetails)

}
