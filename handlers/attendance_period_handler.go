package handlers

import (
	"payroll-system/middleware"
	"payroll-system/models"
	"payroll-system/services"
	"payroll-system/validators"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AttendancePeriodHandler interface {
	CreatePeriod(c *fiber.Ctx) error
}

type AttendancePeriodHandlerModule struct {
	attendancePeriodService services.AttendancePeriodService
}

type AttendancePeriodHandlerOpts struct {
	AttendancePeriodService services.AttendancePeriodService
}

type CreateAttendancePeriodRequest struct {
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

func NewAttendancePeriodHandler(opts *AttendancePeriodHandlerOpts) AttendancePeriodHandler {
	return &AttendancePeriodHandlerModule{attendancePeriodService: opts.AttendancePeriodService}
}

func (h *AttendancePeriodHandlerModule) CreatePeriod(c *fiber.Ctx) error {
	var req CreateAttendancePeriodRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := validators.ValidateStruct(req); err != nil {
		validationError := validators.ParseValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationError,
		})
	}

	employeeIDstr, _ := c.Locals("employee_id").(string)
	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: missing username",
		})
	}
	employeeID, err := uuid.Parse(employeeIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	period := &models.AttendancePeriod{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := h.attendancePeriodService.CreateAttendancePeriod(c.UserContext(), period, employeeID, username, middleware.GetClientIP(c)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(period)
}
