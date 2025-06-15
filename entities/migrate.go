package entities

import (
	"strings"

	"gorm.io/gorm"
)

// AutoMigrate runs all entity migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Route{},
		&Schedule{},
		&Seat{},
		&Booking{},
		&BookingDetail{},
		&Payment{},
	)
}

// CreateIndexes creates additional indexes for better performance
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
