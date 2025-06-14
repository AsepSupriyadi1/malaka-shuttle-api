package entities

import (
	"gorm.io/gorm"
)

type BookingDetail struct {
	gorm.Model
	BookingID     uint    `gorm:"not null;index"`
	SeatID        uint    `gorm:"not null;uniqueIndex"` // One seat can only be in one booking
	PassengerName string  `gorm:"size:100;not null"`
	Price         float64 `gorm:"type:decimal(10,2);not null"`
	// Relations
	Booking Booking `gorm:"foreignKey:BookingID"`
	Seat    Seat    `gorm:"foreignKey:SeatID"`
}
