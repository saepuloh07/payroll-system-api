package seeders

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var firstNames = []string{
	"Ahmad", "Budi", "Siti", "Rina", "Doni", "Lina", "Joko", "Nina", "Fajar", "Dewi",
	"Andi", "Eko", "Mira", "Hendra", "Indra", "Krisna", "Lia", "Mega", "Oki", "Putri",
}

var lastNames = []string{
	"Santoso", "Wibowo", "Prasetyo", "Susanto", "Purnama", "Putri", "Lestari", "Wijaya",
	"Mulyadi", "Kusuma", "Aditya", "Saputra", "Mahendra", "Sinaga", "Pangestu", "Utami",
	"Ramadhan", "Hardiansyah", "Kristianto", "Suryadi",
}

var salaries = []float64{
	4000000, 4500000, 5000000, 5500000,
	6000000, 6500000, 7000000, 7500000,
	8000000, 8500000, 9000000, 10000000,
}

func getRandomSalary() float64 {
	return salaries[rand.Intn(len(salaries))]
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func getRandomName() string {
	first := firstNames[rand.Intn(len(firstNames))]
	last := lastNames[rand.Intn(len(lastNames))]
	return fmt.Sprintf("%s %s", first, last)
}

func SeedEmployees(db *sql.DB) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM public.employees`).Scan(&count)
	if err != nil {
		return
	}

	if count > 0 {
		fmt.Println("Seeding employees already exists")
		return
	}

	adminUsername := "admin"
	adminFullname := "Admin User"
	adminPassword, _ := hashPassword("password123")
	adminRole := "admin"
	adminCreatedAt := time.Now()
	adminUpdatedAt := time.Now()

	// Insert Admin
	_, err = db.Exec(`
        INSERT INTO public.employees (id, username, fullname, password, salary, role, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		uuid.New(), adminUsername, adminFullname, adminPassword, 15000000, adminRole,
		adminCreatedAt, adminUpdatedAt, "system", "system")

	if err != nil {
		log.Fatalf("Failed to insert admin: %v", err)
	}
	fmt.Println("✅ Admin inserted successfully")

	// Insert Employees
	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("employee%d", i)
		fullname := getRandomName()
		password, _ := hashPassword("password123")
		salary := getRandomSalary()
		role := "employee"
		createdAt := time.Now().AddDate(0, 0, -rand.Intn(365))
		updatedAt := createdAt

		_, err := db.Exec(`
            INSERT INTO public.employees (id, username, fullname, password, salary, role, created_at, updated_at, created_by, updated_by)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			uuid.New(), username, fullname, password, salary, role,
			createdAt, updatedAt, "admin", "admin")

		if err != nil {
			log.Printf("Error inserting employee %d: %v", i, err)
			continue
		}
	}

	fmt.Println("🎉 Seeding completed successfully!")
}
