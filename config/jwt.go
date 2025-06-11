package config

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	EmployeeId uuid.UUID `json:"employee_id"`
	Fullname   string    `json:"fullname"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	jwt.RegisteredClaims
}
