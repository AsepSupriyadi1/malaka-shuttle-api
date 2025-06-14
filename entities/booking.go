package entities

import (
	"time"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending             BookingStatus = "pending"
	BookingStatusWaitingVerification BookingStatus = "waiting_verification"
	BookingStatusSuccess             BookingStatus = "success"
	BookingStatusRejected            BookingStatus = "rejected"
	BookingStatusExpired             BookingStatus = "expired"
	BookingStatusCancelled           BookingStatus = "cancelled"
)

type Booking struct {
	gorm.Model
	UserID        uint `gorm:"not null;index"`
	ScheduleID    uint `gorm:"not null;index"`
	BookingTime   time.Time
	Status        BookingStatus `gorm:"type:enum('pending','waiting_verification','success','rejected','expired','cancelled');default:'pending'"`
	ExpiresAt     time.Time     `gorm:"not null"`
	PaymentAmount float64       `gorm:"type:decimal(10,2);not null;default:0"`

	// Relations
	User           User            `gorm:"foreignKey:UserID"`
	Schedule       Schedule        `gorm:"foreignKey:ScheduleID"`
	BookingDetails []BookingDetail `gorm:"foreignKey:BookingID"`
	Payment        *Payment        `gorm:"foreignKey:BookingID"`
}
