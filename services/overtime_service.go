package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"payroll-system/models"
	"payroll-system/repositories"
	"payroll-system/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OvertimeService interface {
	SubmitOvertime(ctx context.Context, overtime *models.Overtime, username, ipAddress string) error
}

type OvertimeServiceModule struct {
	db           *sqlx.DB
	overtimeRepo repositories.OvertimeRepository
	auditLogRepo repositories.AuditLogRepository
}

type OvertimeServiceOpts struct {
	Db           *sqlx.DB
	OvertimeRepo repositories.OvertimeRepository
	AuditLogRepo repositories.AuditLogRepository
}

func NewOvertimeService(opts *OvertimeServiceOpts) OvertimeService {
	return &OvertimeServiceModule{db: opts.Db, overtimeRepo: opts.OvertimeRepo, auditLogRepo: opts.AuditLogRepo}
}

func (s *OvertimeServiceModule) SubmitOvertime(ctx context.Context, overtime *models.Overtime, username, ipAddress string) error {
	now := time.Now()

	workHourStart := time.Date(overtime.Date.Year(), overtime.Date.Month(), overtime.Date.Day(), 9, 0, 0, 0, now.Location())
	workHourEnd := time.Date(overtime.Date.Year(), overtime.Date.Month(), overtime.Date.Day(), 17, 0, 0, 0, now.Location())

	if overtime.Date.Before(workHourStart) || overtime.Date.Equal(workHourEnd) {
		return errors.New("overtime can only be submitted after working hours (after 17:00)")
	}

	if overtime.Hours <= 0 || overtime.Hours > 3 {
		return errors.New("overtime must be between 0.1 and 3 hours")
	}

	existingOvertime, err := s.overtimeRepo.GetTodayOvertime(ctx, overtime.EmployeeID, overtime.Date)
	if err != nil {
		return err
	}
	if existingOvertime != nil {
		newTotal := existingOvertime.Hours + overtime.Hours
		if newTotal > 3 {
			return fmt.Errorf("total overtime cannot exceed 3 hours per day (current total: %.2f)", existingOvertime.Hours)
		}
	}

	overtime.ID = uuid.New()
	overtime.CreatedAt = time.Now()
	overtime.UpdatedAt = time.Now()
	overtime.CreatedBy = username
	overtime.UpdatedBy = username

	newJson, _ := json.Marshal(overtime)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: overtime.EmployeeID,
		TableName:  "overtimes",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   overtime.ID.String(),
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
		return s.overtimeRepo.CreateOvertime(ctx, tx, overtime)
	}); err != nil {
		return err
	}

	return nil
}
