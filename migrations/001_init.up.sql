BEGIN;

CREATE TABLE employees (
    id UUID PRIMARY KEY,
    fullname VARCHAR(150) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(150) NOT NULL,
    salary NUMERIC NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE attendance_periods (
    id UUID PRIMARY KEY,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE attendances (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    fullname TEXT NOT NULL,
    attendance_date TIMESTAMP NOT NULL,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE overtimes (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    fullname TEXT NOT NULL,
    overtime_date TIMESTAMP NOT NULL,
    hours NUMERIC NOT NULL,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE reimbursements (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    fullname TEXT NOT NULL,
    reimbursement_date TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    locked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE payrolls (
    id UUID PRIMARY KEY ,
    attendance_period_id UUID NOT NULL,
    run_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL,
    updated_by VARCHAR(50) NOT NULL
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    employee_id UUID NOT NULL,
    request_id UUID NOT NULL,
    table_name TEXT NOT NULL,
    record_id TEXT NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('create', 'update', 'delete')),
    old_values JSONB,
    new_values JSONB,
    ip_address TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(50) NOT NULL
);

CREATE INDEX idx_audit_log_request_id ON audit_logs(request_id);
CREATE INDEX idx_audit_log_employee_id ON audit_logs(employee_id);
CREATE INDEX idx_audit_log_table_name ON audit_logs(table_name);

COMMIT;