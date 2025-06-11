package handlers

import (
	"payroll-system/middleware"
	"payroll-system/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PayrollHandler interface {
	Run(c *fiber.Ctx) error
}

type PayrollHandlerModule struct {
	payrollService services.PayrollService
}

type PayrollHandlerOpts struct {
	PayrollService services.PayrollService
}

func NewPayrollHandler(opts *PayrollHandlerOpts) PayrollHandler {
	return &PayrollHandlerModule{payrollService: opts.PayrollService}
}

type RunPayrollRequest struct {
	AttendancePeriodID string `json:"attendance_period_id"` // UUID sebagai string
}

func (h *PayrollHandlerModule) Run(c *fiber.Ctx) error {
	var req RunPayrollRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	adminUserIdstr, _ := c.Locals("employee_id").(string)
	adminUser, ok := c.Locals("username").(string)
	if !ok || adminUser == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	periodID, err := uuid.Parse(req.AttendancePeriodID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid attendance period ID format"})
	}

	employeeID, err := uuid.Parse(adminUserIdstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	if err := h.payrollService.RunPayroll(c.UserContext(), periodID, employeeID, adminUser, middleware.GetClientIP(c)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Payroll successfully processed and records are now locked",
	})
}
