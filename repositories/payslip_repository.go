package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"payroll-system/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PayslipRepository interface {
	GetEmployeeData(ctx context.Context, employeeID uuid.UUID) (*models.Payslip, error)
	GetAllEmployeeTakeHomePay(ctx context.Context) ([]models.EmployeeSummary, error)
}

type PayslipRepositoryModule struct {
	db *sqlx.DB
}

type PayslipRepositoryOps struct {
	Db *sqlx.DB
}

func NewPayslipRepository(ops *PayslipRepositoryOps) PayslipRepository {
	return &PayslipRepositoryModule{
		db: ops.Db,
	}
}

func (r *PayslipRepositoryModule) GetEmployeeData(ctx context.Context, employeeID uuid.UUID) (*models.Payslip, error) {
	var emp struct {
		ID       uuid.UUID `db:"id"`
		Username string    `db:"username"`
		Fullname string    `db:"fullname"`
		Salary   float64   `db:"salary"`
	}

	err := r.db.Get(&emp, "SELECT id, username, fullname, salary FROM employees WHERE id = $1", employeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found")
	}

	var attendances []models.AttendanceRecord
	err = r.db.SelectContext(ctx, &attendances, `
        SELECT attendance_date FROM attendances WHERE employee_id = $1`, employeeID)
	if err != nil {
		return nil, err
	}

	var overtimes []models.OvertimeRecord
	err = r.db.SelectContext(ctx, &overtimes, `
        SELECT overtime_date, hours FROM overtimes WHERE employee_id = $1`, employeeID)
	if err != nil {
		return nil, err
	}

	var reims []models.ReimbursementRecord
	err = r.db.SelectContext(ctx, &reims, `
        SELECT description, amount FROM reimbursements WHERE employee_id = $1`, employeeID)
	if err != nil {
		return nil, err
	}

	hourlySalary := math.Round(emp.Salary / (22 * 8))
	takeHome := emp.Salary
	totalOvertimeHour := 0.0
	for _, o := range overtimes {
		totalOvertimeHour += o.Hours
	}
	totalOvertimePay := totalOvertimeHour * hourlySalary
	takeHome += totalOvertimePay

	totalReimbursement := 0.0
	for _, r := range reims {
		totalReimbursement += r.Amount
	}
	takeHome += totalReimbursement

	payslip := &models.Payslip{
		ID:                 uuid.New(),
		EmployeeID:         emp.ID,
		Fullname:           emp.Fullname,
		Salary:             emp.Salary,
		Attendance:         attendances,
		Overtime:           overtimes,
		Reimbursements:     reims,
		TotalOvertimeHour:  totalOvertimeHour,
		TotalOvertimePay:   totalOvertimePay,
		TotalReimbursement: totalReimbursement,
		TakeHomePay:        takeHome,
		CreatedAt:          time.Now(),
		CreatedBy:          emp.Username,
	}

	return payslip, nil
}

func (r *PayslipRepositoryModule) GetAllEmployeeTakeHomePay(ctx context.Context) ([]models.EmployeeSummary, error) {
	var summaries []models.EmployeeSummary

	query := `
        SELECT 
			emp.fullname, 
			emp.salary, 
			COALESCE(SUM(ROUND(emp.salary / (22 * 8)) * ovr.hours),0) AS overtime_pay, 
			COALESCE(SUM(rmb.amount),0) AS reimbursement_pay
		FROM employees emp 
		LEFT JOIN attendances att ON emp.id=att.employee_id  
		LEFT JOIN overtimes ovr ON emp.id=ovr.employee_id 
		LEFT JOIN reimbursements rmb ON emp.id=rmb.employee_id 
		WHERE emp.role = 'employee' Group by emp.fullname, emp.salary;`

	err := r.db.Select(&summaries, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to fetch employee summaries: %v", err)
	}

	return summaries, nil
}
