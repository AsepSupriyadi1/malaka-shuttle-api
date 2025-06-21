package config

import (
	"fmt"
	"log"
	"malakashuttle/entities"
	"math/rand"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database")
	}

	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Route{},
		&entities.Schedule{},
		&entities.Seat{},
		&entities.Booking{},
		&entities.BookingDetail{},
		&entities.Payment{},
	)
}

// SeedData populates the database with sample data for testing
func SeedData(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Create admin user
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := entities.User{
		Email:       "admin@malakashuttle.com",
		Password:    string(adminPassword),
		Role:        "admin",
		FirstName:   "Admin",
		LastName:    "System",
		PhoneNumber: "081234567890",
	}
	if err := db.FirstOrCreate(&admin, entities.User{Email: admin.Email}).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %v", err)
	}

	// Create staff user
	staffPassword, _ := bcrypt.GenerateFromPassword([]byte("staff123"), bcrypt.DefaultCost)
	staff := entities.User{
		Email:       "staff@malakashuttle.com",
		Password:    string(staffPassword),
		Role:        "staff",
		FirstName:   "Staff",
		LastName:    "One",
		PhoneNumber: "081234567891",
	}
	if err := db.FirstOrCreate(&staff, entities.User{Email: staff.Email}).Error; err != nil {
		return fmt.Errorf("failed to create staff user: %v", err)
	}

	// Create 10 regular users
	for i := 1; i <= 10; i++ {
		password, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("user%d123", i)), bcrypt.DefaultCost)
		user := entities.User{
			Email:       fmt.Sprintf("user%d@example.com", i),
			Password:    string(password),
			Role:        "user",
			FirstName:   fmt.Sprintf("User%d", i),
			LastName:    "Test",
			PhoneNumber: fmt.Sprintf("08123456789%d", i),
		}
		if err := db.FirstOrCreate(&user, entities.User{Email: user.Email}).Error; err != nil {
			return fmt.Errorf("failed to create user %d: %v", i, err)
		}
	}
	log.Println("✓ Users seeded successfully")

	// Create routes
	routes := []entities.Route{
		{OriginCity: "Jakarta", DestinationCity: "Bandung"},
		{OriginCity: "Bandung", DestinationCity: "Yogyakarta"},
		{OriginCity: "Yogyakarta", DestinationCity: "Surabaya"},
		{OriginCity: "Surabaya", DestinationCity: "Malang"},
		{OriginCity: "Jakarta", DestinationCity: "Semarang"},
		{OriginCity: "Semarang", DestinationCity: "Solo"},
	}

	for _, route := range routes {
		if err := db.FirstOrCreate(&route, entities.Route{
			OriginCity:      route.OriginCity,
			DestinationCity: route.DestinationCity,
		}).Error; err != nil {
			return fmt.Errorf("failed to create route %s-%s: %v", route.OriginCity, route.DestinationCity, err)
		}
	}
	log.Println("✓ Routes seeded successfully")

	// Get all routes from database
	var dbRoutes []entities.Route
	if err := db.Find(&dbRoutes).Error; err != nil {
		return fmt.Errorf("failed to fetch routes: %v", err)
	}
	// Create 10 schedules total (distribute across routes)
	scheduleCount := 0
	maxSchedules := 10

	for _, route := range dbRoutes {
		if scheduleCount >= maxSchedules {
			break
		}

		// Create 1-2 schedules per route depending on available slots
		schedulesForRoute := 2
		if scheduleCount+schedulesForRoute > maxSchedules {
			schedulesForRoute = maxSchedules - scheduleCount
		}

		for i := 0; i < schedulesForRoute && scheduleCount < maxSchedules; i++ {
			// Random departure time in next 3 days
			dayOffset := rand.Intn(3)
			hourOffset := 8 + rand.Intn(10) // Between 8 AM - 6 PM

			departureTime := time.Now().AddDate(0, 0, dayOffset).Add(time.Duration(hourOffset) * time.Hour)
			arrivalTime := departureTime.Add(time.Duration(2+rand.Intn(3)) * time.Hour) // 2-4 hours journey

			// Random seats between 8-10
			totalSeats := 8 + rand.Intn(3) // 8, 9, or 10 seats

			schedule := entities.Schedule{
				RouteID:        route.ID,
				DepartureTime:  departureTime,
				ArrivalTime:    arrivalTime,
				Price:          50000 + float64(rand.Intn(100000)), // Random price between 50k-150k
				TotalSeats:     totalSeats,
				AvailableSeats: totalSeats,
			}

			if err := db.Create(&schedule).Error; err != nil {
				return fmt.Errorf("failed to create schedule: %v", err)
			}
			scheduleCount++
		}
	}
	log.Println("✓ Schedules seeded successfully")

	// Get all schedules from database
	var dbSchedules []entities.Schedule
	if err := db.Find(&dbSchedules).Error; err != nil {
		return fmt.Errorf("failed to fetch schedules: %v", err)
	}
	// Create seats for each schedule (8-10 seats per schedule)
	for _, schedule := range dbSchedules {
		var seats []entities.Seat
		for seatNum := 1; seatNum <= schedule.TotalSeats; seatNum++ {
			seat := entities.Seat{
				ScheduleID: schedule.ID,
				SeatNumber: fmt.Sprintf("%d", seatNum), // Simple numbering: 1, 2, 3, etc.
				IsBooked:   false,
			}
			seats = append(seats, seat)
		}

		// Create all seats for this schedule
		if err := db.Create(&seats).Error; err != nil {
			return fmt.Errorf("failed to create seats for schedule %d: %v", schedule.ID, err)
		}
	}
	log.Println("✓ Seats seeded successfully")

	log.Println("Database seeding completed successfully!")
	return nil
}

// ResetDatabase drops all tables and recreates them with fresh data
func ResetDatabase(db *gorm.DB) error {
	log.Println("Starting database reset...")

	// Drop all tables
	if err := db.Migrator().DropTable(
		&entities.Payment{},
		&entities.BookingDetail{},
		&entities.Booking{},
		&entities.Seat{},
		&entities.Schedule{},
		&entities.Route{},
		&entities.User{},
	); err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}
	log.Println("✓ All tables dropped")

	// Recreate tables
	if err := AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate tables: %v", err)
	}
	log.Println("✓ Tables recreated")

	// Seed fresh data
	if err := SeedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}

	log.Println("Database reset completed successfully!")
	return nil
}
