package repositories

import (
	"context"
	"database/sql"
	"payroll-system/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ReimbursementRepository interface {
	CreateReimbursement(ctx context.Context, tx *sql.Tx, reimbursement *models.Reimbursement) error
	UpdateReimbursementLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error
}

type ReimbursementRepositoryModule struct {
	db *sqlx.DB
}

type ReimbursementRepositoryOpts struct {
	Db *sqlx.DB
}

func NewReimbursementRepository(opts *ReimbursementRepositoryOpts) ReimbursementRepository {
	return &ReimbursementRepositoryModule{db: opts.Db}
}

func (r *ReimbursementRepositoryModule) CreateReimbursement(ctx context.Context, tx *sql.Tx, rmb *models.Reimbursement) error {
	query := `
    INSERT INTO reimbursements (id, employee_id, fullname, amount, description, reimbursement_date, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := tx.ExecContext(ctx, query,
		rmb.ID,
		rmb.EmployeeID,
		rmb.Fullname,
		rmb.Amount,
		rmb.Description,
		rmb.Date,
		rmb.CreatedAt,
		rmb.UpdatedAt,
		rmb.CreatedBy,
		rmb.UpdatedBy,
	)

	return err
}

func (r *ReimbursementRepositoryModule) UpdateReimbursementLocked(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) error {
	_, err := tx.ExecContext(ctx, `
        UPDATE reimbursements SET locked = TRUE WHERE reimbursement_date BETWEEN $1 AND $2;`,
		startDate, endDate)
	return err
}
