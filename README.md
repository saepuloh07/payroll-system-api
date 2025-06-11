PROJECT NAME: Payroll System API

DESCRIPTION:
A backend API for managing employee attendance, overtime, reimbursements, and payroll processing.
Built with Golang + GoFiber + PostgreSQL.

TECHNOLOGY STACK:
- Programming Language: Golang
- Web Framework: GoFiber
- Database: PostgreSQL
- JWT Authentication
- Input Validation: go-playground/validator
- Middleware: Auth & Role-Based Access Control
- Database Migration: golang-migrate/migrate

DIRECTORY STRUCTURE (Monolithic)
  ```bash
payroll-system/
├── main.go
├── go.mod
├── .env
├── config/
│   ├── database.go
│   └── jwt.go
├── models/
│   ├── employee.go
│   ├── attendance.go
│   ├── overtime.go
│   ├── reimbursement.go
│   ├── attendance_period.go
│   ├── payroll.go
│   ├── payslip.go
│   └── audit_log.go
├── repositories/
│   ├── employee_repository.go
│   ├── attendance_repository.go
│   ├── attendance_period_repository.go
│   ├── overtime_repository.go
│   ├── reimbursement_repository.go
│   ├── payroll_repository.go
│   ├── payslip_repository.go
│   └── audit_log_repository.go
├── services/
│   ├── auth_service.go
│   ├── attendance_service.go
│   ├── attendance_period_service.go
│   ├── overtime_service.go
│   ├── payroll_service.go
│   ├── payslip_service.go
│   └── reimbursement_service.go
├── handlers/
│   ├── auth_handler.go
│   ├── attendance_handler.go
│   ├── attendance_period_handler.go
│   ├── overtime_handler.go
│   ├── payroll_handler.go
│   ├── payslip_handler.go
│   └── reimbursement_handler.go
├── routes/
│   ├── auth_route.go
│   ├── attendance_route.go
│   ├── attendance_period_route.go
│   ├── overtime_route.go
│   ├── payroll_route.go
│   ├── payslip_route.go
│   └── reimbursement_route.go
├── middleware/
│   ├── audit_middleware.go
│   ├── auth_middleware.go
│   └── role_middleware.go
├── validators/
│   ├── validator.go
│   └── errors.go
├── seeders/
│   └── employee_seeder.go
├── migrations/
│   ├── 001_init.down.sql
│   └── 001_init.up.sql
└── README.md
```

DATABASE TABLES

employees
| Column         | Data Type     | Description                      |
|----------------|---------------|----------------------------------|
| id             | UUID          | Primary Key                      |
| username       | TEXT          | Unique username                  |
| fullname       | TEXT          | Full name of the employee        |
| password       | TEXT          | Hashed password                  |
| salary         | NUMERIC       | Monthly base salary              |
| role           | VARCHAR(50)   | 'admin' or 'employee'            |
| created_at     | TIMESTAMP     | Account creation time            |
| updated_at     | TIMESTAMP     | Last update time                 |
| created_by     | TEXT          | Admin who created the record     |
| updated_by     | TEXT          | Admin who last updated the record|

attendances
| Column            | Data Type     | Description                        |
|-------------------|---------------|------------------------------------|
| id                | UUID          | Primary Key                        |
| user_id           | TEXT          | Username                           |
| fullname          | TEXT          | Full name                          |
| attendance_date   | TIMESTAMP     | Date of attendance submission      |
| created_at        | TIMESTAMP     | Time record was created            |
| updated_at        | TIMESTAMP     | Last update time                   |
| created_by        | TEXT          | User who created the record        |
| updated_by        | TEXT          | User who last updated the record   |

overtimes
| Column           | Data Type     | Description                         |
|------------------|---------------|-------------------------------------|
| id               | UUID          | Primary Key                         |
| user_id          | TEXT          | Username                            |
| fullname         | TEXT          | Full name                           |
| overtime_date    | TIMESTAMP     | Date of overtime                    |
| hours            | NUMERIC       | Overtime duration (max 3 hours/day) |
| created_at       | TIMESTAMP     | Record creation time                |
| updated_at       | TIMESTAMP     | Last update time                    |
| created_by       | TEXT          | User who created the record         |
| updated_by       | TEXT          | User who last updated the record    |

reimbursements
| Column         | Data Type     | Description                         |
|----------------|---------------|-------------------------------------|
| id             | UUID          | Primary Key                         |
| user_id        | TEXT          | Username                            |
| fullname       | TEXT          | Full name                           |
| amount         | NUMERIC       | Reimbursement amount                |
| description    | TEXT          | Reason/purpose of reimbursement     |
| created_at     | TIMESTAMP     | Record creation time                |
| updated_at     | TIMESTAMP     | Last update time                    |
| created_by     | TEXT          | User who created the record         |
| updated_by     | TEXT          | User who last updated the record    |

attendance_periods
| Column              | Data Type     | Description                             |
|---------------------|---------------|-----------------------------------------|
| id                  | UUID          | Primary Key                             |
| payroll_id          | INTEGER       | Payroll period ID                       |
| start_date          | TIMESTAMP     | Start date of the period                |
| end_date            | TIMESTAMP     | End date of the period                  |
| created_at          | TIMESTAMP     | Record creation time                    |
| updated_at          | TIMESTAMP     | Record last update time                 |
| created_by          | TEXT          | Admin who created                       |
| updated_by          | TEXT          | Admin who last updated                  |

payrolls
| Column                | Data Type     | Description                             |
|-----------------------|---------------|-----------------------------------------|
| id                    | UUID          | Primary Key                             |
| attendance_period_id  | UUID          | Reference to attendance_periods.id      |
| run_date              | TIMESTAMP     | When payroll was processed              |
| created_by            | TEXT          | Admin who ran the payroll               |
| updated_at            | TIMESTAMP     | Record last update time                 |
| created_by            | TEXT          | Admin who created                       |
| updated_by            | TEXT          | Admin who last updated                  |

audit_logs
| Column             | Data Type     | Description                                |
|---------------------|---------------|-------------------------------------------|
| id                  | UUID          | Primary Key                               |
| request_id          | UUID          | Unique ID per request                     |
| table_name          | TEXT          | Table affected                            |
| record_id           | TEXT          | ID of the affected record                 |
| action              | TEXT          | create/update/delete                      |
| old_values          | TEXT          | Previous values (for updates)             |
| new_values          | TEXT          | New values                                |
| ip_address          | TEXT          | Client IP address                         |
| user_id             | TEXT          | Username of the user                      |
| user_fullname       | TEXT          | Full name of the user                     |
| created_at          | TIMESTAMP     | Timestamp                                 |
| created_by          | TEXT          | User who created                          |

MIDDLEWARE

AuthMiddleware
- Verifies JWT token from Authorization header
- Sets `employee_id`, `username`, `fullname`, and `role` in locals

RequireRole(role string)
- Restricts access based on role (e.g., "admin")

RequestIDMiddleware
- Generates a unique `request_id` for each HTTP request
- Stores it in context and locals for logging/tracing

JWT CLAIMS STRUCTURE

  ```bash
{
  "employee_id": "7cb0a1d0-f03d-4748-bc3f-5e6b41e4f4f9",
  "username": "employee1",
  "fullname": "Budi Santoso",
  "role": "employee",
  "exp": 1719812345
}
```

ENDPOINTS

Auth
| Method | Endpoint       | Access | Description                    |
|--------|----------------|--------|--------------------------------|
| POST   | /register      | Public | Employee registration          |
| POST   | /login         | Public | Login and get JWT Token        |

Attendance
| Method | Endpoint           | Access     | Description                   |
|--------|--------------------|------------|-------------------------------|
| POST   | /attendance/submit | Employee   | Submit daily attendance       |

Overtime
| Method | Endpoint           | Access     | Description                      |
|--------|--------------------|------------|----------------------------------|
| POST   | /overtime/submit   | Employee   | Submit overtime with hour count  |

Reimbursement
| Method | Endpoint               | Access     | Description                          |
|--------|------------------------|------------|--------------------------------------|
| POST   | /reimbursement/submit  | Employee   | Submit reimbursement request         |

Attendance Period (Admin)
| Method | Endpoint                     | Access | Description                              |
|--------|------------------------------|--------|------------------------------------------|
| POST   | /attendance-period           | Admin  | Create new attendance period             |

Payroll (Admin)
| Method | Endpoint     | Access | Description                                  |
|--------|--------------|--------|----------------------------------------------|
| POST   | /payroll/run | Admin  | Run payroll → locks data for that period     |

Payslip (Employee)
| Method | Endpoint           | Access     | Description                          |
|--------|--------------------|------------|--------------------------------------|
| GET    | /payslip/generate  | Employee   | Generate payslip from current data   |

Payslip Summary (Admin)
| Method | Endpoint                   | Access | Description                               |
|--------|----------------------------|--------|-------------------------------------------|
| GET    | /admin/payslip/summary     | Admin  | View summary of all employee payslips     |

CLI COMMANDS

  ```bash
go run main.go db up          -> Migrates database
go run main.go db down        -> Rollback Migrates database
go run main.go seed employees -> Seeds 100 test employees
go run main.go                -> Starts Fiber server
```

ENVIRONMENT CONFIGURATION (.env)

  ```bash
DB_DSN=postgres://postgres:yourpassword@localhost:5432/dbpayroll?sslmode=disable
PORT=3000
JWT_SECRET=rahasia_jwt_ku
```

FEATURES

| Feature                            | Status |
|------------------------------------|--------|
| Register & Login (JWT)             | ✅     |
| CRUD Attendance                    | ✅     |
| Submit Overtime                    | ✅     |
| Submit Reimbursement               | ✅     |
| Attendance Period Management       | ✅     |
| Run Payroll (Admin Only)           | ✅     |
| Generate Payslip (Employee)        | ✅     |
| Summary Payslip (Admin)            | ✅     |
| Seeder & DB Migrations             | ✅     |
| Input Validation                   | ✅     |
| Role-Based Access Control          | ✅     |
| Automated Testing                  | ✅     |
| Audit Logging (IP + request_id)    | ✅     |

AUDITING & LOGGING

- Each request gets a unique `request_id`
- IP address is logged for every significant operation
- Changes to critical records are stored in `audit_logs` table
