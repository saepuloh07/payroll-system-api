package repositories

import (
	"context"
	"database/sql"
	"payroll-system/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendanceRepository interface {
	CreateAttendance(ctx context.Context, tx *sql.Tx, att *models.Attendance) error
	HasAttendanceToday(ctx context.Context, employeeID uuid.UUID, today string) (bool, error)
	UpdateAttendanceLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error
}

type AttendanceRepositoryModule struct {
	db *sqlx.DB
}

type AttendanceRepositoryOpts struct {
	Db *sqlx.DB
}

func NewAttendanceRepository(opts *AttendanceRepositoryOpts) AttendanceRepository {
	return &AttendanceRepositoryModule{db: opts.Db}
}

func (r *AttendanceRepositoryModule) CreateAttendance(ctx context.Context, tx *sql.Tx, att *models.Attendance) error {
	query := `
    INSERT INTO attendances (id, employee_id, fullname, attendance_date, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := tx.ExecContext(ctx, query,
		att.ID,
		att.EmployeeID,
		att.Fullname,
		att.Date,
		att.CreatedAt,
		att.UpdatedAt,
		att.CreatedBy,
		att.UpdatedBy,
	)

	return err
}

func (r *AttendanceRepositoryModule) HasAttendanceToday(ctx context.Context, employeeID uuid.UUID, today string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
        SELECT COUNT(*) FROM attendances
        WHERE employee_id = $1 AND DATE(attendance_date) = DATE($2)`,
		employeeID, today)
	if count > 0 {
		return true, nil
	}
	return false, err
}

func (r *AttendanceRepositoryModule) UpdateAttendanceLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE attendances SET locked = TRUE WHERE attendance_date BETWEEN $1 AND $2`,
		startDate, endDate)
	return err
}
