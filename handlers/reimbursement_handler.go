package handlers

import (
	"payroll-system/middleware"
	"payroll-system/models"
	"payroll-system/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReimbursementHandler interface {
	Submit(c *fiber.Ctx) error
}

type ReimbursementHandlerModule struct {
	reimbursementService services.ReimbursementService
}

type ReimbursementHandlerOpts struct {
	ReimbursementService services.ReimbursementService
}

func NewReimbursementHandler(opts *ReimbursementHandlerOpts) ReimbursementHandler {
	return &ReimbursementHandlerModule{reimbursementService: opts.ReimbursementService}
}

type SubmitReimbursementRequest struct {
	Date        time.Time `json:"reimbursement_date"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
}

func (h *ReimbursementHandlerModule) Submit(c *fiber.Ctx) error {
	var req SubmitReimbursementRequest

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

	rmb := &models.Reimbursement{
		EmployeeID:  employeeID,
		Fullname:    fullname,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	}

	if err := h.reimbursementService.SubmitReimbursement(c.UserContext(), rmb, username, middleware.GetClientIP(c)); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":     "Reimbursement submitted successfully",
		"fullname":    fullname,
		"amount":      req.Amount,
		"description": req.Description,
	})
}
