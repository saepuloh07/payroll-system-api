package routes

import (
	"payroll-system/handlers"
	"payroll-system/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPayslipRoute(app *fiber.App, handler handlers.PayslipHandler) {
	route := app.Group("/payslip")
	route.Use(middleware.AuthMiddleware, middleware.RequireRole("employee"))

	route.Get("/generate", handler.Generate)

	adminRoute := app.Group("/admin/payslip")
	adminRoute.Use(middleware.AuthMiddleware, middleware.RequireRole("admin"))

	adminRoute.Get("/summary", handler.GeneratePayslipSummary)
}
