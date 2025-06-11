package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAttendancePeriodRoute(app *fiber.App, handler handlers.AttendancePeriodHandler) {
	route := app.Group("/attendance-period")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("admin"))

	route.Post("/", handler.CreatePeriod)
}
