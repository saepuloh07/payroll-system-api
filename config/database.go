package config

import (
	"fmt"
	"log"
	"os"
	"payroll-system/seeders"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")
	// dbUser := os.Getenv("DB_USER")
	// dbPass := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbSslMode := os.Getenv("DB_SSLMODE")

	// Construct the full DSN URL (Connection String)
	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSslMode)
	// dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbName, dbSslMode)

	// Print DSN for debugging (optional)
	fmt.Println("Connecting to DB with DSN:", dsn)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db

	fmt.Println("Connected to the database")
}

func RunMigrations() {
	migrationPath := "file://migrations"
	dsn := os.Getenv("DB_DSN")

	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	if len(os.Args) < 3 {
		fmt.Println("Missing subcommand for 'db'. Usage: db [up|down]")
		return
	}

	switch os.Args[2] {
	case "up":
		fmt.Println("🚀 Running database migrations...")
		if err := m.Up(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to apply")
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			fmt.Println("✅ Migrations applied successfully")
		}

	case "down":
		fmt.Println("Rolling back migrations...")
		if err := m.Down(); err != nil {
			if err.Error() == "no change" {
				fmt.Println("No migrations to apply")
			} else {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			fmt.Println("✅ Migrations rolled back successfully")
		}

	default:
		fmt.Println("Unknown subcommand:", os.Args[2])
	}
}

func RunSeeders() {
	dbCon := DB.DB
	switch os.Args[2] {
	case "employees":
		seeders.SeedEmployees(dbCon)
	case "all":
		seeders.SeedEmployees(dbCon)
		// Tambahkan seeder lain di sini
	default:
		fmt.Printf("Seeder '%s' not found\n", os.Args[2])
	}
}
