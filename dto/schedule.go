package dto

import (
	"fmt"
	"malakashuttle/entities"
	"time"
)

// CreateScheduleRequest - DTO untuk request create schedule (Admin)
type CreateScheduleRequest struct {
	RouteID       uint    `json:"route_id" validate:"required" binding:"required"`
	DepartureTime string  `json:"departure_time" validate:"required" binding:"required"` // Format: "YYYY-MM-DD HH:mm"
	ArrivalTime   string  `json:"arrival_time" validate:"required" binding:"required"`   // Format: "YYYY-MM-DD HH:mm"
	Price         float64 `json:"price" validate:"required,gt=0" binding:"required"`
	TotalSeats    int     `json:"total_seats" validate:"required,gt=0" binding:"required"`
}

// UpdateScheduleRequest - DTO untuk request update schedule (Admin)
// Note: TotalSeats tidak bisa diubah untuk menghindari konflik data
type UpdateScheduleRequest struct {
	RouteID       *uint    `json:"route_id,omitempty"`
	DepartureTime *string  `json:"departure_time,omitempty"` // Format: "YYYY-MM-DD HH:mm"
	ArrivalTime   *string  `json:"arrival_time,omitempty"`   // Format: "YYYY-MM-DD HH:mm"
	Price         *float64 `json:"price,omitempty"`
}

// ScheduleSearchRequest - DTO untuk pencarian schedule (User)
type ScheduleSearchRequest struct {
	Origin        string `form:"origin" validate:"required" binding:"required"`
	Destination   string `form:"destination" validate:"required" binding:"required"`
	DepartureDate string `form:"departure_date" validate:"required" binding:"required"` // Format: "YYYY-MM-DD"
	Page          int    `form:"page,default=1"`
	Limit         int    `form:"limit,default=10"`
}

// ScheduleResponse - DTO untuk response schedule (unified untuk admin dan user)
type ScheduleResponse struct {
	ID             uint       `json:"id"`
	Origin         string     `json:"origin"`
	Destination    string     `json:"destination"`
	DepartureTime  string     `json:"departure_time"` // Format: "YYYY-MM-DD HH:mm"
	ArrivalTime    string     `json:"arrival_time"`   // Format: "YYYY-MM-DD HH:mm"
	Price          float64    `json:"price"`
	TotalSeats     int        `json:"total_seats,omitempty"` // Bisa null untuk user
	AvailableSeats int        `json:"available_seats"`
	Duration       string     `json:"duration"`
	CreatedAt      *time.Time `json:"created_at,omitempty"` // Bisa null untuk user
	UpdatedAt      *time.Time `json:"updated_at,omitempty"` // Bisa null untuk user
}

// ScheduleListResponse - DTO untuk response list schedule dengan pagination
type ScheduleListResponse struct {
	Results    []ScheduleResponse `json:"results"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
	TotalCount int64              `json:"total_count"`
}

// ToScheduleResponse - Convert entity to response DTO
func ToScheduleResponse(schedule entities.Schedule, includeAdminFields bool) ScheduleResponse {
	// Load timezone Indonesia (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60) // Fallback ke WIB +7
	}

	// Convert to WIB timezone and format to "YYYY-MM-DD HH:mm"
	departureTimeWIB := schedule.DepartureTime.In(loc)
	arrivalTimeWIB := schedule.ArrivalTime.In(loc)

	departureTimeStr := departureTimeWIB.Format("2006-01-02 15:04")
	arrivalTimeStr := arrivalTimeWIB.Format("2006-01-02 15:04")

	// Calculate duration using WIB times
	duration := arrivalTimeWIB.Sub(departureTimeWIB)
	durationStr := fmt.Sprintf("%dh %dm", int(duration.Hours()), int(duration.Minutes())%60)

	response := ScheduleResponse{
		ID:             schedule.ID,
		Origin:         schedule.Route.OriginCity,
		Destination:    schedule.Route.DestinationCity,
		DepartureTime:  departureTimeStr,
		ArrivalTime:    arrivalTimeStr,
		Price:          schedule.Price,
		AvailableSeats: schedule.AvailableSeats,
		Duration:       durationStr,
	}

	// Include admin-only fields if requested
	if includeAdminFields {
		response.TotalSeats = schedule.TotalSeats
		response.CreatedAt = &schedule.CreatedAt
		response.UpdatedAt = &schedule.UpdatedAt
	}

	return response
}

// ToScheduleListResponse - Convert entity list to response DTO with pagination
func ToScheduleListResponse(schedules []entities.Schedule, page, limit int, totalCount int64, includeAdminFields bool) ScheduleListResponse {
	var scheduleResponses []ScheduleResponse
	for _, schedule := range schedules {
		scheduleResponses = append(scheduleResponses, ToScheduleResponse(schedule, includeAdminFields))
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	return ScheduleListResponse{
		Results:    scheduleResponses,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}
