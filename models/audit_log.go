package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	EmployeeID uuid.UUID `json:"employee_id" db:"employee_id"`
	RequestID  uuid.UUID `json:"request_id" db:"request_id"`
	TableName  string    `json:"table_name" db:"table_name"`
	RecordID   string    `json:"record_id" db:"record_id"`
	Action     string    `json:"action" db:"action"` // create, update, delete
	OldValues  string    `json:"old_values,omitempty" db:"old_values"`
	NewValues  string    `json:"new_values,omitempty" db:"new_values"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
}
