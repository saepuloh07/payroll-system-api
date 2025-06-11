package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"payroll-system/config"
	"payroll-system/models"
	"payroll-system/repositories"
	"payroll-system/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterEmployee(ctx context.Context, emp *models.Employee, ipAddress string) error
	Login(ctx context.Context, username, password string) (string, error)
}

type AuthServiceModule struct {
	db           *sqlx.DB
	employeeRepo repositories.EmployeeRepository
	auditLogRepo repositories.AuditLogRepository
}

type AuthServiceOpts struct {
	Db           *sqlx.DB
	EmployeeRepo repositories.EmployeeRepository
	AuditLogRepo repositories.AuditLogRepository
}

func NewAuthService(opts *AuthServiceOpts) AuthService {
	return &AuthServiceModule{db: opts.Db, employeeRepo: opts.EmployeeRepo, auditLogRepo: opts.AuditLogRepo}
}

func (s *AuthServiceModule) RegisterEmployee(ctx context.Context, emp *models.Employee, ipAddress string) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.DefaultCost)
	emp.Password = string(hashedPassword)
	emp.ID = uuid.New()
	emp.CreatedAt = time.Now()
	emp.UpdatedAt = time.Now()
	emp.CreatedBy = "SYSTEM"
	emp.UpdatedBy = "SYSTEM"

	newJson, _ := json.Marshal(emp)
	jsonString := string(newJson)
	logModel := &models.AuditLog{
		ID:         uuid.New(),
		EmployeeID: emp.ID,
		TableName:  "attendances",
		RequestID:  ctx.Value("request_id").(uuid.UUID),
		RecordID:   emp.ID.String(),
		OldValues:  "",
		NewValues:  jsonString,
		Action:     "create",
		IPAddress:  ipAddress,
		CreatedAt:  time.Now(),
		CreatedBy:  "SYSTEM",
	}

	if err := utils.WrapTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := s.auditLogRepo.CreateLog(ctx, tx, logModel); err != nil {
			return err
		}
		return s.employeeRepo.CreateEmployee(ctx, tx, emp)
	}); err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceModule) Login(ctx context.Context, username, password string) (string, error) {
	emp, err := s.employeeRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(emp.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &config.Claims{
		EmployeeId: emp.ID,
		Fullname:   emp.Fullname,
		Username:   emp.Username,
		Role:       emp.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
