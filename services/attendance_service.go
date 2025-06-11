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

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AttendanceService interface {
	SubmitAttendance(ctx context.Context, attendance *models.Attendance, username, ipAddress string) error
}

type AttendanceServiceModule struct {
	db             *sqlx.DB
	attendanceRepo repositories.AttendanceRepository
	auditLogRepo   repositories.AuditLogRepository
}

type AttendanceServiceOpts struct {
	Db             *sqlx.DB
	AttendanceRepo repositories.AttendanceRepository
	AuditLogRepo   repositories.AuditLogRepository
}

func NewAttendanceService(opts *AttendanceServiceOpts) AttendanceService {
	return &AttendanceServiceModule{db: opts.Db, attendanceRepo: opts.AttendanceRepo, auditLogRepo: opts.AuditLogRepo}
}

func (s *AttendanceServiceModule) SubmitAttendance(ctx context.Context, attendance *models.Attendance, username, ipAddress string) error {
	now := time.Now()
	dateOnly := time.Date(attendance.Date.Year(), attendance.Date.Month(), attendance.Date.Day(), 0, 0, 0, 0, now.Location())

	// Cek apakah sudah submit hari ini
	hasToday, err := s.attendanceRepo.HasAttendanceToday(ctx, attendance.EmployeeID, dateOnly.Format("2006-01-02"))
	if err != nil {
		return err
	}
	if hasToday {
		return errors.New("already submitted attendance for today")
	}

	// Cek apakah weekend
	weekday := now.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return errors.New("cannot submit attendance on weekends")
	}

	attendance.ID = uuid.New()
	attendance.Date = dateOnly
	attendance.CreatedAt = now
	attendance.UpdatedAt = now
	attendance.CreatedBy = username
	attendance.UpdatedBy = username

	newJson, _ := json.Marshal(attendance)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: attendance.EmployeeID,
		TableName:  "attendances",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   attendance.ID.String(),
		OldValues:  "",
		NewValues:  jsonString,
		Action:     "create",
		IPAddress:  ipAddress,
		CreatedAt:  now,
		CreatedBy:  username,
	}

	if err := utils.WrapTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.auditLogRepo.CreateLog(ctx, tx, logModel); err != nil {
			return err
		}
		return s.attendanceRepo.CreateAttendance(ctx, tx, attendance)
	}); err != nil {
		return err
	}

	return nil
}
