package handlers

import (
	"payroll-system/middleware"
	"payroll-system/models"
	"payroll-system/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OvertimeHandler interface {
	Submit(c *fiber.Ctx) error
}

type OvertimeHandlerModule struct {
	overtimeService services.OvertimeService
}

type OvertimeHandlerOpts struct {
	OvertimeService services.OvertimeService
}

func NewOvertimeHandler(opts *OvertimeHandlerOpts) OvertimeHandler {
	return &OvertimeHandlerModule{overtimeService: opts.OvertimeService}
}

type SubmitOvertimeRequest struct {
	Date  time.Time `json:"overtime_date"` // format: "2025-04-05T00:00:00Z"
	Hours float64   `json:"hours"`
}

func (h *OvertimeHandlerModule) Submit(c *fiber.Ctx) error {
	var req SubmitOvertimeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	username, ok := c.Locals("username").(string)
	employeeIDstr, _ := c.Locals("employee_id").(string)
	fullname, _ := c.Locals("fullname").(string)

	if !ok || username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	employeeID, err := uuid.Parse(employeeIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	overtime := &models.Overtime{
		EmployeeID: employeeID,
		Fullname:   fullname,
		Date:       req.Date,
		Hours:      req.Hours,
	}

	if err := h.overtimeService.SubmitOvertime(c.UserContext(), overtime, username, middleware.GetClientIP(c)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":  "Overtime submitted successfully",
		"fullname": fullname,
		"hours":    req.Hours,
		"date":     req.Date,
	})
}
