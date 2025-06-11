package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupOvertimeRoute(app *fiber.App, handler handlers.OvertimeHandler) {
	route := app.Group("/overtime")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("employee"))

	route.Post("/submit", handler.Submit)
}
