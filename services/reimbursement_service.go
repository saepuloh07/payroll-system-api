package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"payroll-system/models"
	"payroll-system/repositories"
	"payroll-system/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReimbursementService interface {
	SubmitReimbursement(ctx context.Context, rmb *models.Reimbursement, username, ipAddress string) error
}

type ReimbursementServiceModule struct {
	db                *sqlx.DB
	reimbursementRepo repositories.ReimbursementRepository
	auditLogRepo      repositories.AuditLogRepository
}

type ReimbursementServiceOpts struct {
	Db                *sqlx.DB
	ReimbursementRepo repositories.ReimbursementRepository
	AuditLogRepo      repositories.AuditLogRepository
}

func NewReimbursementService(opts *ReimbursementServiceOpts) ReimbursementService {
	return &ReimbursementServiceModule{db: opts.Db, reimbursementRepo: opts.ReimbursementRepo, auditLogRepo: opts.AuditLogRepo}
}

func (s *ReimbursementServiceModule) SubmitReimbursement(ctx context.Context, rmb *models.Reimbursement, username, ipAddress string) error {
	if rmb.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	rmb.ID = uuid.New()
	rmb.CreatedAt = time.Now()
	rmb.UpdatedAt = time.Now()
	rmb.CreatedBy = username
	rmb.UpdatedBy = username

	newJson, _ := json.Marshal(rmb)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: rmb.EmployeeID,
		TableName:  "reimbursements",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   rmb.ID.String(),
		OldValues:  "",
		NewValues:  jsonString,
		Action:     "create",
		IPAddress:  ipAddress,
		CreatedAt:  time.Now(),
		CreatedBy:  username,
	}

	if err := utils.WrapTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.auditLogRepo.CreateLog(ctx, tx, logModel); err != nil {
			return err
		}
		return s.reimbursementRepo.CreateReimbursement(ctx, tx, rmb)
	}); err != nil {
		return err
	}

	return nil
}
