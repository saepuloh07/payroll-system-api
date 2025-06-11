package models

import (
	"time"

	"github.com/google/uuid"
)

type Reimbursement struct {
	ID          uuid.UUID `json:"id" db:"id"`
	EmployeeID  uuid.UUID `json:"employee_id" db:"employee_id"`
	Fullname    string    `json:"fullname" db:"fullname"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"reimbursement_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	UpdatedBy   string    `json:"updated_by" db:"updated_by"`
}
