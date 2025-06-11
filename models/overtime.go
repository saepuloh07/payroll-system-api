package models

import (
	"time"

	"github.com/google/uuid"
)

type Overtime struct {
	ID         uuid.UUID `json:"id" db:"id"`
	EmployeeID uuid.UUID `json:"employee_id" db:"employee_id"`
	Fullname   string    `json:"fullname" db:"fullname"`
	Date       time.Time `json:"date" db:"overtime_date"`
	Hours      float64   `json:"hours" db:"hours"`
	Locked     bool      `json:"locked" db:"locked"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	UpdatedBy  string    `json:"updated_by" db:"updated_by"`
}

type OvertimeHourCount struct {
	EmployeeID uuid.UUID `json:"employee_id" db:"employee_id"`
	Hours      float64   `json:"hours" db:"hours"`
}
