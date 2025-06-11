package repositories

import (
	"context"
	"database/sql"
	"errors"
	"payroll-system/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OvertimeRepository interface {
	CreateOvertime(ctx context.Context, tx *sql.Tx, ot *models.Overtime) error
	GetTodayOvertime(ctx context.Context, employeeID uuid.UUID, date time.Time) (*models.OvertimeHourCount, error)
	UpdateOvertimeLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error
}

type OvertimeRepositoryModule struct {
	db *sqlx.DB
}

type OvertimeRepositoryOpts struct {
	Db *sqlx.DB
}

func NewOvertimeRepository(opts *OvertimeRepositoryOpts) OvertimeRepository {
	return &OvertimeRepositoryModule{db: opts.Db}
}

func (r *OvertimeRepositoryModule) CreateOvertime(ctx context.Context, tx *sql.Tx, ot *models.Overtime) error {
	query := `
    INSERT INTO overtimes (id, employee_id, fullname, overtime_date, hours, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := tx.ExecContext(ctx, query,
		ot.ID,
		ot.EmployeeID,
		ot.Fullname,
		ot.Date,
		ot.Hours,
		ot.CreatedAt,
		ot.UpdatedAt,
		ot.CreatedBy,
		ot.UpdatedBy,
	)

	return err
}

func (r *OvertimeRepositoryModule) GetTodayOvertime(ctx context.Context, employeeID uuid.UUID, date time.Time) (*models.OvertimeHourCount, error) {
	var ot models.OvertimeHourCount
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	err := r.db.GetContext(ctx, &ot, `
        SELECT employee_id, COALESCE(SUM(hours),0) as hours FROM overtimes 
        WHERE employee_id = $1 AND overtime_date BETWEEN $2 AND $3
        GROUP BY employee_id`,
		employeeID, startOfDay, endOfDay)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error(err.Error())
		return nil, err
	}

	return &ot, nil
}

func (r *OvertimeRepositoryModule) UpdateOvertimeLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE overtimes SET locked = TRUE WHERE overtime_date BETWEEN $1 AND $2;`,
		startDate, endDate)
	return err
}
