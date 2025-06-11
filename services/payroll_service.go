package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"payroll-system/models"
	"payroll-system/repositories"
	"payroll-system/utils"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PayrollService interface {
	RunPayroll(ctx context.Context, periodID uuid.UUID, adminID uuid.UUID, adminUsername, ipAddress string) error
}

type PayrollServiceModule struct {
	db                *sqlx.DB
	payrollRepo       repositories.PayrollRepository
	attPeriodRepo     repositories.AttendancePeriodRepository
	attendanceRepo    repositories.AttendanceRepository
	overtimeRepo      repositories.OvertimeRepository
	reimbursementRepo repositories.ReimbursementRepository
	auditLogRepo      repositories.AuditLogRepository
}

type PayrollServiceOpts struct {
	Db                *sqlx.DB
	PayrollRepo       repositories.PayrollRepository
	AttPeriodRepo     repositories.AttendancePeriodRepository
	AttendanceRepo    repositories.AttendanceRepository
	OvertimeRepo      repositories.OvertimeRepository
	ReimbursementRepo repositories.ReimbursementRepository
	AuditLogRepo      repositories.AuditLogRepository
}

func NewPayrollService(opts *PayrollServiceOpts) PayrollService {
	return &PayrollServiceModule{
		db:                opts.Db,
		payrollRepo:       opts.PayrollRepo,
		attPeriodRepo:     opts.AttPeriodRepo,
		attendanceRepo:    opts.AttendanceRepo,
		overtimeRepo:      opts.OvertimeRepo,
		reimbursementRepo: opts.ReimbursementRepo,
		auditLogRepo:      opts.AuditLogRepo,
	}
}

func (s *PayrollServiceModule) RunPayroll(ctx context.Context, periodID uuid.UUID, adminID uuid.UUID, adminUsername, ipAddress string) error {
	attendancePeriod, err := s.attPeriodRepo.GetAttendancePeriodByAttendancePeriodID(ctx, periodID)
	if err != nil {
		return err
	}

	if attendancePeriod == nil {
		return errors.New("attendance period not found")
	}

	alreadyRun, err := s.payrollRepo.IsPayrollRun(ctx, periodID)
	if err != nil {
		return err
	}
	if alreadyRun {
		return errors.New("payroll for this period has already been processed")
	}

	payroll := &models.Payroll{
		ID:                 uuid.New(),
		AttendancePeriodID: periodID,
		RunDate:            time.Now(),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          adminUsername,
		UpdatedBy:          adminUsername,
	}

	newJson, _ := json.Marshal(payroll)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: adminID,
		TableName:  "payrolls",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   payroll.ID.String(),
		OldValues:  "",
		NewValues:  jsonString,
		Action:     "create",
		IPAddress:  ipAddress,
		CreatedAt:  time.Now(),
		CreatedBy:  adminUsername,
	}

	if err := utils.WrapTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.auditLogRepo.CreateLog(ctx, tx, logModel); err != nil {
			return err
		}

		if err := s.payrollRepo.CreatePayrollRun(ctx, tx, payroll); err != nil {
			return err
		}

		if err := s.attendanceRepo.UpdateAttendanceLocked(ctx, tx, attendancePeriod.StartDate, attendancePeriod.EndDate); err != nil {
			log.Error(err)
			return err
		}

		if err := s.overtimeRepo.UpdateOvertimeLocked(ctx, tx, attendancePeriod.StartDate, attendancePeriod.EndDate); err != nil {
			log.Error(err)
			return err
		}

		if err := s.reimbursementRepo.UpdateReimbursementLocked(ctx, tx, attendancePeriod.StartDate, attendancePeriod.EndDate); err != nil {
			log.Error(err)
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
