package models

import (
	"time"

	"github.com/google/uuid"
)

type Payroll struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	AttendancePeriodID uuid.UUID `json:"attendance_period_id" db:"attendance_period_id"`
	RunDate            time.Time `json:"run_date" db:"run_date"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	UpdatedBy          string    `json:"updated_by" db:"updated_by"`
}
