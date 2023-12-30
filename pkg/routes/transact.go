package routes

import (
	"github.com/aikintech/scim-api/pkg/controllers"
	"github.com/aikintech/scim-api/pkg/middlewares"
	"github.com/gofiber/fiber/v2"
)

func MountTransactionRoutes(app *fiber.App) {
	route := app.Group("/transactions")

	transactionController := controllers.NewTransactionController()
	jwtAuthWare := middlewares.JWTMiddleware("access")

	route.Get("/", jwtAuthWare, transactionController.GetTransactions)
	route.Post("/", jwtAuthWare, transactionController.Transact)
	route.Get("/:transactionId", jwtAuthWare, transactionController.GetTransaction)

	// Webhooks
	route.Post("/webhooks/paystack", transactionController.PaystackWebhook)
}
