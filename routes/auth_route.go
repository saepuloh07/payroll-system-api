package routes

import (
	"payroll-system/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(app *fiber.App, handler handlers.AuthHandler) {
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)
}
