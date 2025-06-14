package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// BookingStatus enum
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusPaid      BookingStatus = "paid"
	BookingStatusExpired   BookingStatus = "expired"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking model untuk business logic layer
type Booking struct {
	ID          uint          `json:"id"`
	UserID      uint          `json:"user_id"`
	ScheduleID  uint          `json:"schedule_id"`
	BookingTime time.Time     `json:"booking_time"`
	Status      BookingStatus `json:"status"`
	ExpiresAt   time.Time     `json:"expires_at"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// BookingWithDetails model untuk operasi yang memerlukan detail lengkap
type BookingWithDetails struct {
	Booking
	User           *User            `json:"user,omitempty"`
	Schedule       *Schedule        `json:"schedule,omitempty"`
	BookingDetails []*BookingDetail `json:"booking_details,omitempty"`
	Payment        *Payment         `json:"payment,omitempty"`
}

// BookingHandler untuk converter operations
type BookingHandler struct{}

// NewBookingHandler membuat instance baru dari BookingHandler
func NewBookingHandler() *BookingHandler {
	return &BookingHandler{}
}

// FromEntity mengkonversi entity.Booking ke model.Booking
func (h *BookingHandler) FromEntity(entity *entities.Booking) *Booking {
	if entity == nil {
		return nil
	}

	return &Booking{
		ID:          entity.ID,
		UserID:      entity.UserID,
		ScheduleID:  entity.ScheduleID,
		BookingTime: entity.BookingTime,
		Status:      BookingStatus(entity.Status),
		ExpiresAt:   entity.ExpiresAt,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.Booking ke entity.Booking
func (h *BookingHandler) ToEntity(model *Booking) *entities.Booking {
	if model == nil {
		return nil
	}

	return &entities.Booking{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		UserID:      model.UserID,
		ScheduleID:  model.ScheduleID,
		BookingTime: model.BookingTime,
		Status:      entities.BookingStatus(model.Status),
		ExpiresAt:   model.ExpiresAt,
	}
}

// FromEntityWithDetails mengkonversi entity.Booking ke model.BookingWithDetails
func (h *BookingHandler) FromEntityWithDetails(
	entity *entities.Booking,
	userHandler *UserHandler,
	scheduleHandler *ScheduleHandler,
	bookingDetailHandler *BookingDetailHandler,
	paymentHandler *PaymentHandler,
) *BookingWithDetails {
	if entity == nil {
		return nil
	}

	booking := &BookingWithDetails{
		Booking: Booking{
			ID:          entity.ID,
			UserID:      entity.UserID,
			ScheduleID:  entity.ScheduleID,
			BookingTime: entity.BookingTime,
			Status:      BookingStatus(entity.Status),
			ExpiresAt:   entity.ExpiresAt,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		},
	}

	// Convert user jika ada
	if userHandler != nil {
		booking.User = userHandler.FromEntity(&entity.User)
	}

	// Convert schedule jika ada
	if scheduleHandler != nil {
		booking.Schedule = scheduleHandler.FromEntity(&entity.Schedule)
	}

	// Convert booking details jika ada
	if len(entity.BookingDetails) > 0 && bookingDetailHandler != nil {
		details := make([]*BookingDetail, len(entity.BookingDetails))
		for i, detail := range entity.BookingDetails {
			details[i] = bookingDetailHandler.FromEntity(&detail)
		}
		booking.BookingDetails = details
	}

	// Convert payment jika ada
	if entity.Payment != nil && paymentHandler != nil {
		booking.Payment = paymentHandler.FromEntity(entity.Payment)
	}

	return booking
}

// FromEntityList mengkonversi slice entity.Booking ke slice model.Booking
func (h *BookingHandler) FromEntityList(entities []*entities.Booking) []*Booking {
	if entities == nil {
		return nil
	}

	models := make([]*Booking, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var BookingHandlerInstance = NewBookingHandler()
