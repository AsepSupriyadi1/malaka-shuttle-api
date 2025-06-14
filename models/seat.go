package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// Seat model untuk business logic layer
type Seat struct {
	ID         uint      `json:"id"`
	ScheduleID uint      `json:"schedule_id"`
	SeatNumber string    `json:"seat_number"`
	IsBooked   bool      `json:"is_booked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// SeatHandler untuk converter operations
type SeatHandler struct{}

// NewSeatHandler membuat instance baru dari SeatHandler
func NewSeatHandler() *SeatHandler {
	return &SeatHandler{}
}

// FromEntity mengkonversi entity.Seat ke model.Seat
func (h *SeatHandler) FromEntity(entity *entities.Seat) *Seat {
	if entity == nil {
		return nil
	}

	return &Seat{
		ID:         entity.ID,
		ScheduleID: entity.ScheduleID,
		SeatNumber: entity.SeatNumber,
		IsBooked:   entity.IsBooked,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.Seat ke entity.Seat
func (h *SeatHandler) ToEntity(model *Seat) *entities.Seat {
	if model == nil {
		return nil
	}

	return &entities.Seat{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		ScheduleID: model.ScheduleID,
		SeatNumber: model.SeatNumber,
		IsBooked:   model.IsBooked,
	}
}

// FromEntityList mengkonversi slice entity.Seat ke slice model.Seat
func (h *SeatHandler) FromEntityList(entities []*entities.Seat) []*Seat {
	if entities == nil {
		return nil
	}

	models := make([]*Seat, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var SeatHandlerInstance = NewSeatHandler()
