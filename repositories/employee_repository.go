package repositories

import (
	"context"
	"database/sql"
	"payroll-system/models"

	"github.com/jmoiron/sqlx"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, tx *sql.Tx, emp *models.Employee) error
	GetEmployeeByUsername(ctx context.Context, username string) (*models.Employee, error)
}

type EmployeeRepositoryModule struct {
	db *sqlx.DB
}

type EmployeeRepositoryOpts struct {
	Db *sqlx.DB
}

func NewEmployeeRepository(opts *EmployeeRepositoryOpts) EmployeeRepository {
	return &EmployeeRepositoryModule{db: opts.Db}
}

func (r *EmployeeRepositoryModule) CreateEmployee(ctx context.Context, tx *sql.Tx, emp *models.Employee) error {
	query := `
    INSERT INTO employees (id, fullname, username, password, salary, role, created_at, updated_at, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id`

	err := tx.QueryRowContext(ctx, query,
		emp.ID,
		emp.Fullname,
		emp.Username,
		emp.Password,
		emp.Salary,
		emp.Role,
		emp.CreatedAt,
		emp.UpdatedAt,
		emp.CreatedBy,
		emp.UpdatedBy,
	).Scan(&emp.ID)

	return err
}

func (r *EmployeeRepositoryModule) GetEmployeeByUsername(ctx context.Context, username string) (*models.Employee, error) {
	var emp models.Employee
	err := r.db.GetContext(ctx, &emp, "SELECT * FROM employees WHERE username = $1", username)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}
