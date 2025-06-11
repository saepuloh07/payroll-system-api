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

type AttendancePeriodRepository interface {
	CreatePeriod(ctx context.Context, tx *sql.Tx, period *models.AttendancePeriod) error
	GetAttendancePeriodByAttendancePeriodID(ctx context.Context, periodID uuid.UUID) (*models.AttendancePeriod, error)
	CheckAvailableStartAndEndDate(ctx context.Context, startDate, endDate time.Time) (bool, error)
}

type AttendancePeriodRepositoryModule struct {
	db *sqlx.DB
}

type AttendancePeriodRepositoryOpts struct {
	Db *sqlx.DB
}

func NewAttendancePeriodRepository(opts *AttendancePeriodRepositoryOpts) AttendancePeriodRepository {
	return &AttendancePeriodRepositoryModule{db: opts.Db}
}

func (r *AttendancePeriodRepositoryModule) CreatePeriod(ctx context.Context, tx *sql.Tx, period *models.AttendancePeriod) error {
	query := `
    INSERT INTO attendance_periods (id, start_date, end_date, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id`

	err := tx.QueryRowContext(ctx, query,
		period.ID,
		period.StartDate,
		period.EndDate,
		period.CreatedAt,
		period.UpdatedAt,
		period.CreatedBy,
		period.UpdatedBy,
	).Scan(&period.ID)

	return err
}

func (r *AttendancePeriodRepositoryModule) GetAttendancePeriodByAttendancePeriodID(ctx context.Context, periodID uuid.UUID) (*models.AttendancePeriod, error) {
	var emp models.AttendancePeriod
	err := r.db.GetContext(ctx, &emp, "SELECT * FROM attendance_periods WHERE id = $1", periodID)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *AttendancePeriodRepositoryModule) CheckAvailableStartAndEndDate(ctx context.Context, startDate, endDate time.Time) (bool, error) {
	var emp models.AttendancePeriod
	err := r.db.GetContext(ctx, &emp, "SELECT * FROM attendance_periods WHERE ($1 >= start_date AND $1 <= end_date) OR ($2 >= start_date AND $2 <= end_date) limit 1", startDate, endDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		log.Error(err)
		return false, err
	}
	return true, nil
}
