package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPayrollRoute(app *fiber.App, handler handlers.PayrollHandler) {
	route := app.Group("/payroll")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("admin"))

	route.Post("/run", handler.Run)
}
