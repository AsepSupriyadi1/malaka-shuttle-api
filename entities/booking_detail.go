package entities

import (
	"gorm.io/gorm"
)

type BookingDetail struct {
	gorm.Model
	BookingID     uint    `gorm:"not null;index"`
	SeatID        uint    `gorm:"not null;index"` // Changed from uniqueIndex to regular index
	PassengerName string  `gorm:"size:100;not null"`
	Price         float64 `gorm:"type:decimal(10,2);not null"`
	// Relations
	Booking Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE"`
	Seat    Seat    `gorm:"foreignKey:SeatID;constraint:OnDelete:CASCADE"`
}
