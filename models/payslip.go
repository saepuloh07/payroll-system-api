package models

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceRecord struct {
	Date time.Time `json:"date" db:"attendance_date"`
}

type OvertimeRecord struct {
	Date  time.Time `json:"date" db:"overtime_date"`
	Hours float64   `json:"hours" db:"hours"`
}

type ReimbursementRecord struct {
	Description string  `json:"description" db:"description"`
	Amount      float64 `json:"amount" db:"amount"`
}

type Payslip struct {
	ID                 uuid.UUID             `json:"id" db:"id"`
	EmployeeID         uuid.UUID             `json:"employee_id" db:"employee_id"`
	Fullname           string                `json:"fullname" db:"fullname"`
	Salary             float64               `json:"salary" db:"salary"`
	Attendance         []AttendanceRecord    `json:"attendance"`
	Overtime           []OvertimeRecord      `json:"overtime"`
	Reimbursements     []ReimbursementRecord `json:"reimbursements"`
	TotalOvertimeHour  float64               `json:"total_overtime_hour"`
	TotalOvertimePay   float64               `json:"total_overtime_pay"`
	TotalReimbursement float64               `json:"total_reimbursement"`
	TakeHomePay        float64               `json:"take_home_pay"`
	CreatedAt          time.Time             `json:"created_at" db:"created_at"`
	CreatedBy          string                `json:"created_by" db:"created_by"`
}

type EmployeeSummary struct {
	Fullname         string  `json:"fullname" db:"fullname"`
	Salary           float64 `json:"salary" db:"salary"`
	OvertimePay      float64 `json:"overtime_pay" db:"overtime_pay"`
	ReimbursementPay float64 `json:"reimbursement_pay" db:"reimbursement_pay"`
	TakeHomePay      float64 `json:"take_home_pay"`
}

type PayslipSummary struct {
	EmployeeSummaries []EmployeeSummary `json:"employee_summaries"`
	TotalCompanyPay   float64           `json:"total_company_pay"`
}
