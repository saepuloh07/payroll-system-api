package repositories

import (
	"context"
	"database/sql"
	"payroll-system/models"

	"github.com/jmoiron/sqlx"
)

type AuditLogRepository interface {
	CreateLog(ctx context.Context, tx *sql.Tx, log *models.AuditLog) error
}

type AuditLogRepositoryModule struct {
	db *sqlx.DB
}

type AuditLogRepositoryOpts struct {
	Db *sqlx.DB
}

func NewAuditLogRepository(opts *AuditLogRepositoryOpts) AuditLogRepository {
	return &AuditLogRepositoryModule{db: opts.Db}
}

func (r *AuditLogRepositoryModule) CreateLog(ctx context.Context, tx *sql.Tx, log *models.AuditLog) error {
	query := `
    INSERT INTO audit_logs (id, employee_id, request_id, table_name, record_id, action, old_values, new_values, ip_address, created_at, created_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := tx.ExecContext(ctx, query,
		log.ID,
		log.EmployeeID,
		log.RequestID,
		log.TableName,
		log.RecordID,
		log.Action,
		log.OldValues,
		log.NewValues,
		log.IPAddress,
		log.CreatedAt,
		log.CreatedBy,
	)

	return err
}
