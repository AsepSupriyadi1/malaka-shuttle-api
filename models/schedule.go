package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// Schedule model untuk business logic layer
type Schedule struct {
	ID             uint      `json:"id"`
	RouteID        uint      `json:"route_id"`
	DepartureTime  time.Time `json:"departure_time"`
	ArrivalTime    time.Time `json:"arrival_time"`
	Price          float64   `json:"price"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ScheduleWithDetails model untuk operasi yang memerlukan detail lengkap
type ScheduleWithDetails struct {
	Schedule
	Route *Route  `json:"route,omitempty"`
	Seats []*Seat `json:"seats,omitempty"`
}

// ScheduleHandler untuk converter operations
type ScheduleHandler struct{}

// NewScheduleHandler membuat instance baru dari ScheduleHandler
func NewScheduleHandler() *ScheduleHandler {
	return &ScheduleHandler{}
}

// FromEntity mengkonversi entity.Schedule ke model.Schedule
func (h *ScheduleHandler) FromEntity(entity *entities.Schedule) *Schedule {
	if entity == nil {
		return nil
	}

	return &Schedule{
		ID:             entity.ID,
		RouteID:        entity.RouteID,
		DepartureTime:  entity.DepartureTime,
		ArrivalTime:    entity.ArrivalTime,
		Price:          entity.Price,
		TotalSeats:     entity.TotalSeats,
		AvailableSeats: entity.AvailableSeats,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.Schedule ke entity.Schedule
func (h *ScheduleHandler) ToEntity(model *Schedule) *entities.Schedule {
	if model == nil {
		return nil
	}

	return &entities.Schedule{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		RouteID:        model.RouteID,
		DepartureTime:  model.DepartureTime,
		ArrivalTime:    model.ArrivalTime,
		Price:          model.Price,
		TotalSeats:     model.TotalSeats,
		AvailableSeats: model.AvailableSeats,
	}
}

// FromEntityWithDetails mengkonversi entity.Schedule ke model.ScheduleWithDetails
func (h *ScheduleHandler) FromEntityWithDetails(entity *entities.Schedule, routeHandler *RouteHandler, seatHandler *SeatHandler) *ScheduleWithDetails {
	if entity == nil {
		return nil
	}

	schedule := &ScheduleWithDetails{
		Schedule: Schedule{
			ID:             entity.ID,
			RouteID:        entity.RouteID,
			DepartureTime:  entity.DepartureTime,
			ArrivalTime:    entity.ArrivalTime,
			Price:          entity.Price,
			TotalSeats:     entity.TotalSeats,
			AvailableSeats: entity.AvailableSeats,
			CreatedAt:      entity.CreatedAt,
			UpdatedAt:      entity.UpdatedAt,
		},
	}

	// Convert route jika ada
	if routeHandler != nil {
		schedule.Route = routeHandler.FromEntity(&entity.Route)
	}

	// Convert seats jika ada
	if len(entity.Seats) > 0 && seatHandler != nil {
		seats := make([]*Seat, len(entity.Seats))
		for i, seat := range entity.Seats {
			seats[i] = seatHandler.FromEntity(&seat)
		}
		schedule.Seats = seats
	}

	return schedule
}

// FromEntityList mengkonversi slice entity.Schedule ke slice model.Schedule
func (h *ScheduleHandler) FromEntityList(entities []*entities.Schedule) []*Schedule {
	if entities == nil {
		return nil
	}

	models := make([]*Schedule, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var ScheduleHandlerInstance = NewScheduleHandler()
