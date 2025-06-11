package main

import (
	"context"
	"os"
	"os/signal"
	"payroll-system/config"
	"payroll-system/handlers"
	"payroll-system/middleware"
	"payroll-system/repositories"
	"payroll-system/routes"
	"payroll-system/services"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitDB()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "db":
			config.RunMigrations()
		case "seed":
			config.RunSeeders()
		}
		return
	}

	app := fiber.New()

	app.Use(middleware.RequestIDMiddleware)

	logRepo := repositories.NewAuditLogRepository(&repositories.AuditLogRepositoryOpts{Db: config.DB})

	empRepo := repositories.NewEmployeeRepository(&repositories.EmployeeRepositoryOpts{Db: config.DB})
	authService := services.NewAuthService(&services.AuthServiceOpts{Db: config.DB, EmployeeRepo: empRepo, AuditLogRepo: logRepo})
	authHandler := handlers.NewAuthHandler(&handlers.AuthHandlerOpts{AuthService: authService})

	attPeriodRepo := repositories.NewAttendancePeriodRepository(&repositories.AttendancePeriodRepositoryOpts{Db: config.DB})
	attPeriodService := services.NewAttendancePeriodService(&services.AttendancePeriodServiceOpts{Db: config.DB, AttendancePeriodRepo: attPeriodRepo, AuditLogRepo: logRepo})
	attPeriodHandler := handlers.NewAttendancePeriodHandler(&handlers.AttendancePeriodHandlerOpts{AttendancePeriodService: attPeriodService})

	attendanceRepo := repositories.NewAttendanceRepository(&repositories.AttendanceRepositoryOpts{Db: config.DB})
	attendanceService := services.NewAttendanceService(&services.AttendanceServiceOpts{Db: config.DB, AttendanceRepo: attendanceRepo, AuditLogRepo: logRepo})
	attendanceHandler := handlers.NewAttendanceHandler(&handlers.AttendanceHandlerOpts{AttendanceService: attendanceService})

	overtimeRepo := repositories.NewOvertimeRepository(&repositories.OvertimeRepositoryOpts{Db: config.DB})
	overtimeService := services.NewOvertimeService(&services.OvertimeServiceOpts{Db: config.DB, OvertimeRepo: overtimeRepo, AuditLogRepo: logRepo})
	overtimeHandler := handlers.NewOvertimeHandler(&handlers.OvertimeHandlerOpts{OvertimeService: overtimeService})

	rmbRepo := repositories.NewReimbursementRepository(&repositories.ReimbursementRepositoryOpts{Db: config.DB})
	rmbService := services.NewReimbursementService(&services.ReimbursementServiceOpts{Db: config.DB, ReimbursementRepo: rmbRepo, AuditLogRepo: logRepo})
	rmbHandler := handlers.NewReimbursementHandler(&handlers.ReimbursementHandlerOpts{ReimbursementService: rmbService})

	payrollRepo := repositories.NewPayrollRepository(&repositories.PayrollRepositoryOpts{Db: config.DB})
	payrollService := services.NewPayrollService(&services.PayrollServiceOpts{
		Db:                config.DB,
		PayrollRepo:       payrollRepo,
		AttPeriodRepo:     attPeriodRepo,
		AttendanceRepo:    attendanceRepo,
		OvertimeRepo:      overtimeRepo,
		ReimbursementRepo: rmbRepo,
		AuditLogRepo:      logRepo,
	})
	payrollHandler := handlers.NewPayrollHandler(&handlers.PayrollHandlerOpts{PayrollService: payrollService})

	payslipRepo := repositories.NewPayslipRepository(&repositories.PayslipRepositoryOps{Db: config.DB})
	payslipService := services.NewPayslipService(&services.PayslipServiceOpts{PayslipRepo: payslipRepo})
	payslipHandler := handlers.NewPayslipHandler(&handlers.PayslipHandlerOpts{PayslipService: payslipService})

	routes.SetupAuthRoute(app, authHandler)
	routes.SetupAttendancePeriodRoute(app, attPeriodHandler)
	routes.SetupAttendanceRoute(app, attendanceHandler)
	routes.SetupReimbursementRoute(app, rmbHandler)
	routes.SetupOvertimeRoute(app, overtimeHandler)
	routes.SetupPayrollRoute(app, payrollHandler)
	routes.SetupPayslipRoute(app, payslipHandler)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		panic(err)
	}

	println("Server shutdown gracefully")
}
