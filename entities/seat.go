package entities

import (
	"gorm.io/gorm"
)

type Seat struct {
	gorm.Model
	ScheduleID uint   `gorm:"not null;index"`
	SeatNumber string `gorm:"size:10;not null"`
	IsBooked   bool   `gorm:"type:boolean;default:false;not null"`

	// Relations
	Schedule      Schedule       `gorm:"foreignKey:ScheduleID"`
	BookingDetail *BookingDetail `gorm:"foreignKey:SeatID"`
}

// TableName returns the table name for Seat
func (Seat) TableName() string {
	return "seats"
}
