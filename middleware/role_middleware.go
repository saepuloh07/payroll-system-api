package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// RequireRole memastikan pengguna memiliki role tertentu
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Access denied. Required role: %s", requiredRole),
			})
		}
		return c.Next()
	}
}
