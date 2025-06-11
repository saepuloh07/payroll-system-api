package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"payroll-system/models"
	"payroll-system/repositories"
	"payroll-system/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendancePeriodService interface {
	CreateAttendancePeriod(ctx context.Context, period *models.AttendancePeriod, employeeID uuid.UUID, username, ipAddress string) error
}

type AttendancePeriodServiceModule struct {
	db                   *sqlx.DB
	attendancePeriodRepo repositories.AttendancePeriodRepository
	auditLogRepo         repositories.AuditLogRepository
}

type AttendancePeriodServiceOpts struct {
	Db                   *sqlx.DB
	AttendancePeriodRepo repositories.AttendancePeriodRepository
	AuditLogRepo         repositories.AuditLogRepository
}

func NewAttendancePeriodService(opts *AttendancePeriodServiceOpts) AttendancePeriodService {
	return &AttendancePeriodServiceModule{db: opts.Db, attendancePeriodRepo: opts.AttendancePeriodRepo, auditLogRepo: opts.AuditLogRepo}
}

func (s *AttendancePeriodServiceModule) CreateAttendancePeriod(ctx context.Context, period *models.AttendancePeriod, employeeID uuid.UUID, username, ipAddress string) error {
	period.ID = uuid.New()
	period.CreatedAt = time.Now()
	period.UpdatedAt = time.Now()
	period.CreatedBy = username
	period.UpdatedBy = username

	newJson, _ := json.Marshal(period)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		TableName:  "attendance_periods",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   period.ID.String(),
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
		return s.attendancePeriodRepo.CreatePeriod(ctx, tx, period)
	}); err != nil {
		return err
	}

	return nil
}
