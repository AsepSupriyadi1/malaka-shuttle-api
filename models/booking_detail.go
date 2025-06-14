package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// BookingDetail model untuk business logic layer
type BookingDetail struct {
	ID            uint      `json:"id"`
	BookingID     uint      `json:"booking_id"`
	SeatID        uint      `json:"seat_id"`
	PassengerName string    `json:"passenger_name"`
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BookingDetailWithSeat model untuk operasi yang memerlukan detail seat
type BookingDetailWithSeat struct {
	BookingDetail
	Seat *Seat `json:"seat,omitempty"`
}

// BookingDetailHandler untuk converter operations
type BookingDetailHandler struct{}

// NewBookingDetailHandler membuat instance baru dari BookingDetailHandler
func NewBookingDetailHandler() *BookingDetailHandler {
	return &BookingDetailHandler{}
}

// FromEntity mengkonversi entity.BookingDetail ke model.BookingDetail
func (h *BookingDetailHandler) FromEntity(entity *entities.BookingDetail) *BookingDetail {
	if entity == nil {
		return nil
	}

	return &BookingDetail{
		ID:            entity.ID,
		BookingID:     entity.BookingID,
		SeatID:        entity.SeatID,
		PassengerName: entity.PassengerName,
		Price:         entity.Price,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.BookingDetail ke entity.BookingDetail
func (h *BookingDetailHandler) ToEntity(model *BookingDetail) *entities.BookingDetail {
	if model == nil {
		return nil
	}

	return &entities.BookingDetail{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		BookingID:     model.BookingID,
		SeatID:        model.SeatID,
		PassengerName: model.PassengerName,
		Price:         model.Price,
	}
}

// FromEntityWithSeat mengkonversi entity.BookingDetail ke model.BookingDetailWithSeat
func (h *BookingDetailHandler) FromEntityWithSeat(entity *entities.BookingDetail, seatHandler *SeatHandler) *BookingDetailWithSeat {
	if entity == nil {
		return nil
	}

	detail := &BookingDetailWithSeat{
		BookingDetail: BookingDetail{
			ID:            entity.ID,
			BookingID:     entity.BookingID,
			SeatID:        entity.SeatID,
			PassengerName: entity.PassengerName,
			Price:         entity.Price,
			CreatedAt:     entity.CreatedAt,
			UpdatedAt:     entity.UpdatedAt,
		},
	}

	// Convert seat jika ada
	if seatHandler != nil {
		detail.Seat = seatHandler.FromEntity(&entity.Seat)
	}

	return detail
}

// FromEntityList mengkonversi slice entity.BookingDetail ke slice model.BookingDetail
func (h *BookingDetailHandler) FromEntityList(entities []*entities.BookingDetail) []*BookingDetail {
	if entities == nil {
		return nil
	}

	models := make([]*BookingDetail, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var BookingDetailHandlerInstance = NewBookingDetailHandler()
