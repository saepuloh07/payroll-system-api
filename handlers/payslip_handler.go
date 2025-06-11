package handlers

import (
	"payroll-system/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PayslipHandler interface {
	Generate(c *fiber.Ctx) error
	GeneratePayslipSummary(c *fiber.Ctx) error
}

type PayslipHandlerModule struct {
	payslipService services.PayslipService
}
type PayslipHandlerOpts struct {
	PayslipService services.PayslipService
}

func NewPayslipHandler(opts *PayslipHandlerOpts) PayslipHandler {
	return &PayslipHandlerModule{
		payslipService: opts.PayslipService,
	}
}

func (h *PayslipHandlerModule) Generate(c *fiber.Ctx) error {
	employeeIDstr, ok := c.Locals("employee_id").(string)
	if !ok || employeeIDstr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	employeeID, err := uuid.Parse(employeeIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	payslip, err := h.payslipService.GeneratePayslip(c.UserContext(), employeeID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(payslip)
}

func (h *PayslipHandlerModule) GeneratePayslipSummary(c *fiber.Ctx) error {
	summary, err := h.payslipService.GeneratePayslipSummary(c.UserContext())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(summary)
}
