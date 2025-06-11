package handlers

import (
	"payroll-system/middleware"
	"payroll-system/models"
	"payroll-system/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AttendanceHandler interface {
	Submit(c *fiber.Ctx) error
}

type AttendanceHandlerModule struct {
	attendanceService services.AttendanceService
}

type AttendanceHandlerOpts struct {
	AttendanceService services.AttendanceService
}

type SubmitAttendanceRequest struct {
	Date time.Time `json:"attendance_date"` // format: "2025-04-05T00:00:00Z"
}

func NewAttendanceHandler(opts *AttendanceHandlerOpts) AttendanceHandler {
	return &AttendanceHandlerModule{attendanceService: opts.AttendanceService}
}

func (h *AttendanceHandlerModule) Submit(c *fiber.Ctx) error {
	var req SubmitAttendanceRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	employeeIDstr, ok := c.Locals("employee_id").(string)
	fullname, _ := c.Locals("fullname").(string)
	username, _ := c.Locals("username").(string)

	if !ok || employeeIDstr == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	employeeID, err := uuid.Parse(employeeIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	attendance := &models.Attendance{
		EmployeeID: employeeID,
		Fullname:   fullname,
	}

	if err := h.attendanceService.SubmitAttendance(c.UserContext(), attendance, username, middleware.GetClientIP(c)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Attendance submitted successfully"})
}
