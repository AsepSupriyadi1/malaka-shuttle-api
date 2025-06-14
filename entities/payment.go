package entities

import (
	"time"

	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	gorm.Model
	BookingID     uint          `gorm:"not null;uniqueIndex"` // One booking = one payment
	PaymentMethod string        `gorm:"size:50;not null"`
	PaymentStatus PaymentStatus `gorm:"type:enum('pending','success','failed');default:'pending'"`
	PaymentDate   *time.Time    `gorm:"null"`
	ProofImageURL string        `gorm:"type:text"`

	// Relations
	Booking Booking `gorm:"foreignKey:BookingID"`
}
