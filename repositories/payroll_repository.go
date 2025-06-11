package repositories

import (
	"context"
	"database/sql"
	"payroll-system/models"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PayrollRepository interface {
	IsPayrollRun(ctx context.Context, periodID uuid.UUID) (bool, error)
	CreatePayrollRun(ctx context.Context, tx *sql.Tx, payroll *models.Payroll) error
}

type PayrollRepositoryModule struct {
	db *sqlx.DB
}

type PayrollRepositoryOpts struct {
	Db *sqlx.DB
}

func NewPayrollRepository(opts *PayrollRepositoryOpts) PayrollRepository {
	return &PayrollRepositoryModule{db: opts.Db}
}

func (r *PayrollRepositoryModule) IsPayrollRun(ctx context.Context, periodID uuid.UUID) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
        SELECT COUNT(*) FROM payrolls WHERE attendance_period_id = $1`, periodID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err.Error())
		return false, err
	}
	return count > 0, nil
}

func (r *PayrollRepositoryModule) CreatePayrollRun(ctx context.Context, tx *sql.Tx, payroll *models.Payroll) error {
	query := `
    INSERT INTO payrolls (id, attendance_period_id, run_date, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := tx.ExecContext(ctx, query,
		payroll.ID,
		payroll.AttendancePeriodID,
		payroll.RunDate,
		payroll.CreatedAt,
		payroll.UpdatedAt,
		payroll.CreatedBy,
		payroll.UpdatedBy,
	)
	return err
}
