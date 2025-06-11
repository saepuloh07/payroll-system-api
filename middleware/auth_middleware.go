package middleware

import (
	"payroll-system/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware memeriksa token JWT pada header Authorization
func AuthMiddleware(c *fiber.Ctx) error {
	// Ambil header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Format: Bearer <token>
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token format",
		})
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &config.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(*config.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Simpan informasi user ke context (opsional)
	c.Locals("employee_id", claims.EmployeeId.String())
	c.Locals("fullname", claims.Fullname)
	c.Locals("username", claims.Username)
	c.Locals("role", claims.Role)

	return c.Next()
}
