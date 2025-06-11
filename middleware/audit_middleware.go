package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware menambahkan request_id ke context
func RequestIDMiddleware(c *fiber.Ctx) error {
	requestID := uuid.New()
	c.Locals("request_id", requestID)

	ctx := context.WithValue(c.UserContext(), "request_id", requestID)
	c.SetUserContext(ctx)

	return c.Next()
}

// GetClientIP mendeteksi IP client
func GetClientIP(c *fiber.Ctx) string {
	ip := c.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.IP()
	}
	return ip
}
