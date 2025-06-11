package handlers

import (
	"payroll-system/middleware"
	"payroll-system/models"
	"payroll-system/requests"
	"payroll-system/services"
	"payroll-system/validators"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
}

type AuthHandlerModule struct {
	authService services.AuthService
}

type AuthHandlerOpts struct {
	AuthService services.AuthService
}

func NewAuthHandler(opts *AuthHandlerOpts) AuthHandler {
	return &AuthHandlerModule{authService: opts.AuthService}
}

func (h *AuthHandlerModule) Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		validationError := validators.ParseValidationErrors(err)
		return c.Status(400).JSON(validationError)
	}

	emp := models.Employee{
		Fullname: req.Fullname,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		Salary:   req.Salary,
	}

	if err := h.authService.RegisterEmployee(c.UserContext(), &emp, middleware.GetClientIP(c)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(emp)
}

func (h *AuthHandlerModule) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.authService.Login(c.UserContext(), input.Username, input.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"token": token})
}
