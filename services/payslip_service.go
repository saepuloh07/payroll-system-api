package services

import (
	"context"
	"payroll-system/models"
	"payroll-system/repositories"

	"github.com/google/uuid"
)

type PayslipService interface {
	GeneratePayslip(ctx context.Context, employeeID uuid.UUID) (*models.Payslip, error)
	GeneratePayslipSummary(ctx context.Context) (*models.PayslipSummary, error)
}

type PayslipServiceModule struct {
	payslipRepo repositories.PayslipRepository
}

type PayslipServiceOpts struct {
	PayslipRepo repositories.PayslipRepository
}

func NewPayslipService(opts *PayslipServiceOpts) PayslipService {
	return &PayslipServiceModule{
		payslipRepo: opts.PayslipRepo,
	}
}

func (s *PayslipServiceModule) GeneratePayslip(ctx context.Context, employeeID uuid.UUID) (*models.Payslip, error) {
	return s.payslipRepo.GetEmployeeData(ctx, employeeID)
}

func (s *PayslipServiceModule) GeneratePayslipSummary(ctx context.Context) (*models.PayslipSummary, error) {
	employeeSummaries, err := s.payslipRepo.GetAllEmployeeTakeHomePay(ctx)
	if err != nil {
		return nil, err
	}

	totalCompanyPay := 0.0
	for i := 0; i < len(employeeSummaries); i++ {
		employeeSummaries[i].TakeHomePay = employeeSummaries[i].Salary + employeeSummaries[i].OvertimePay + employeeSummaries[i].ReimbursementPay
		totalCompanyPay += employeeSummaries[i].TakeHomePay
	}

	summary := &models.PayslipSummary{
		EmployeeSummaries: employeeSummaries,
		TotalCompanyPay:   totalCompanyPay,
	}

	return summary, nil
}
