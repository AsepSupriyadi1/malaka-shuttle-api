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
func CreateIndexes(db *gorm.DB) error { // Hapus IF NOT EXISTS - biarkan error jika index sudah ada
	if err := db.Exec("CREATE UNIQUE INDEX idx_seats_schedule_seat ON seats(schedule_id, seat_number)").Error; err != nil {
		// Log error tapi jangan return error jika index sudah ada
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	if err := db.Exec("CREATE INDEX idx_schedules_departure ON schedules(departure_time)").Error; err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	if err := db.Exec("CREATE INDEX idx_bookings_status ON bookings(status)").Error; err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	if err := db.Exec("CREATE INDEX idx_bookings_expires ON bookings(expires_at)").Error; err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") {
			return err
		}
	}

	return nil
}
