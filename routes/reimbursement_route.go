package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupReimbursementRoute(app *fiber.App, handler handlers.ReimbursementHandler) {
	route := app.Group("/reimbursement")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("employee"))

	route.Post("/submit", handler.Submit)
}
