package models

import (
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

// Route model untuk business logic layer
type Route struct {
	ID              uint      `json:"id"`
	OriginCity      string    `json:"origin_city"`
	DestinationCity string    `json:"destination_city"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// RouteWithSchedules model untuk operasi yang memerlukan schedule
type RouteWithSchedules struct {
	Route
	Schedules []*Schedule `json:"schedules,omitempty"`
}

// RouteHandler untuk converter operations
type RouteHandler struct{}

// NewRouteHandler membuat instance baru dari RouteHandler
func NewRouteHandler() *RouteHandler {
	return &RouteHandler{}
}

// FromEntity mengkonversi entity.Route ke model.Route
func (h *RouteHandler) FromEntity(entity *entities.Route) *Route {
	if entity == nil {
		return nil
	}

	return &Route{
		ID:              entity.ID,
		OriginCity:      entity.OriginCity,
		DestinationCity: entity.DestinationCity,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}

// ToEntity mengkonversi model.Route ke entity.Route
func (h *RouteHandler) ToEntity(model *Route) *entities.Route {
	if model == nil {
		return nil
	}

	return &entities.Route{
		Model: gorm.Model{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		OriginCity:      model.OriginCity,
		DestinationCity: model.DestinationCity,
	}
}

// FromEntityWithSchedules mengkonversi entity.Route ke model.RouteWithSchedules
func (h *RouteHandler) FromEntityWithSchedules(entity *entities.Route, scheduleHandler *ScheduleHandler) *RouteWithSchedules {
	if entity == nil {
		return nil
	}

	route := &RouteWithSchedules{
		Route: Route{
			ID:              entity.ID,
			OriginCity:      entity.OriginCity,
			DestinationCity: entity.DestinationCity,
			CreatedAt:       entity.CreatedAt,
			UpdatedAt:       entity.UpdatedAt,
		},
	}

	// Convert schedules jika ada
	if len(entity.Schedules) > 0 && scheduleHandler != nil {
		schedules := make([]*Schedule, len(entity.Schedules))
		for i, schedule := range entity.Schedules {
			schedules[i] = scheduleHandler.FromEntity(&schedule)
		}
		route.Schedules = schedules
	}

	return route
}

// FromEntityList mengkonversi slice entity.Route ke slice model.Route
func (h *RouteHandler) FromEntityList(entities []*entities.Route) []*Route {
	if entities == nil {
		return nil
	}

	models := make([]*Route, len(entities))
	for i, entity := range entities {
		models[i] = h.FromEntity(entity)
	}
	return models
}

// Global handler instance
var RouteHandlerInstance = NewRouteHandler()
