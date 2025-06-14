package entities

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	RouteID        uint      `gorm:"not null;index"`
	DepartureTime  time.Time `gorm:"not null"`
	ArrivalTime    time.Time `gorm:"not null"`
	Price          float64   `gorm:"type:decimal(10,2);not null"`
	TotalSeats     int       `gorm:"not null"`
	AvailableSeats int       `gorm:"not null"`
	// Relations
	Route    Route     `gorm:"foreignKey:RouteID"`
	Seats    []Seat    `gorm:"foreignKey:ScheduleID"`
	Bookings []Booking `gorm:"foreignKey:ScheduleID"`
}
