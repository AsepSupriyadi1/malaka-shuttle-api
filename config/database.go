package config

import (
	"fmt"
	"malakashuttle/entities"
	"os"
	"strings"

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

func CreateIndexes(db *gorm.DB) error {
	indexes := []string{
		"CREATE UNIQUE INDEX idx_seats_schedule_seat ON seats(schedule_id, seat_number)",
		"CREATE INDEX idx_schedules_departure ON schedules(departure_time)",
		"CREATE INDEX idx_bookings_status ON bookings(status)",
		"CREATE INDEX idx_bookings_expires ON bookings(expires_at)",
	}

	for _, index := range indexes {
		if err := db.Exec(index).Error; err != nil {
			// Ignore error if index already exists
			if !strings.Contains(strings.ToLower(err.Error()), "already exists") &&
				!strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				return err
			}
		}
	}

	// MySQL doesn't support partial indexes with WHERE clause
	// Create a regular unique index on seat_id instead
	partialIndex := `CREATE UNIQUE INDEX idx_booking_details_active_seat ON booking_details(seat_id)`

	if err := db.Exec(partialIndex).Error; err != nil {
		// Ignore error if index already exists
		if !strings.Contains(strings.ToLower(err.Error()), "already exists") &&
			!strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return err
		}
	}

	return nil
}
