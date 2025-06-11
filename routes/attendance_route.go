package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAttendanceRoute(app *fiber.App, handler handlers.AttendanceHandler) {
	route := app.Group("/attendance")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("employee"))

	route.Post("/submit", handler.Submit)
}
