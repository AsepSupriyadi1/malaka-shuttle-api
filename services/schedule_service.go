package services

import (
	"errors"
	"fmt"
	"malakashuttle/dto"
	"malakashuttle/entities"
	"malakashuttle/repositories"
	"time"
)

type ScheduleService struct {
	scheduleRepo *repositories.ScheduleRepository
}

func NewScheduleService(scheduleRepo *repositories.ScheduleRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
	}
}

// CreateSchedule - Create new schedule (Admin only)
func (s *ScheduleService) CreateSchedule(req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	// Load timezone Indonesia (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60) // Fallback ke WIB +7
	}

	// Validasi format waktu dengan timezone lokal
	departureTime, err := time.ParseInLocation("2006-01-02 15:04", req.DepartureTime, loc)
	if err != nil {
		return nil, errors.New("invalid departure_time format, use YYYY-MM-DD HH:mm")
	}

	arrivalTime, err := time.ParseInLocation("2006-01-02 15:04", req.ArrivalTime, loc)
	if err != nil {
		return nil, errors.New("invalid arrival_time format, use YYYY-MM-DD HH:mm")
	}
	// Validasi waktu logis
	if arrivalTime.Before(departureTime) || arrivalTime.Equal(departureTime) {
		return nil, errors.New("arrival_time must be after departure_time")
	}

	// Validasi waktu tidak boleh di masa lalu (bandingkan dengan waktu sekarang di timezone yang sama)
	nowInLoc := time.Now().In(loc)
	if departureTime.Before(nowInLoc) {
		return nil, errors.New("departure_time cannot be in the past")
	}

	// Cek apakah route exists
	routeExists, err := s.scheduleRepo.CheckRouteExists(req.RouteID)
	if err != nil {
		return nil, fmt.Errorf("error checking route: %v", err)
	}
	if !routeExists {
		return nil, errors.New("route not found")
	}

	// Validasi business rules
	if req.TotalSeats <= 0 {
		return nil, errors.New("total_seats must be greater than 0")
	}
	if req.TotalSeats > 50 { // Batasi max 50 seats per schedule
		return nil, errors.New("total_seats cannot exceed 50")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	// Create schedule entity
	schedule := entities.Schedule{
		RouteID:        req.RouteID,
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          req.Price,
		TotalSeats:     req.TotalSeats,
		AvailableSeats: req.TotalSeats, // Available seats sama dengan total seats saat create
	}

	// Save to database (dengan transaction untuk create seats juga)
	if err := s.scheduleRepo.CreateSchedule(&schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %v", err)
	}
	// Get created schedule dengan route relation untuk response
	createdSchedule, err := s.scheduleRepo.GetScheduleByID(schedule.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created schedule: %v", err)
	}

	response := dto.ToScheduleResponse(*createdSchedule, true) // true untuk admin fields
	return &response, nil
}

// UpdateSchedule - Update schedule (Admin only)
func (s *ScheduleService) UpdateSchedule(id uint, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	// Cek apakah schedule exists
	existingSchedule, err := s.scheduleRepo.GetScheduleByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}
	updates := make(map[string]interface{})

	// Load timezone Indonesia (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60) // Fallback ke WIB +7
	}

	// Validasi dan set departure_time
	if req.DepartureTime != nil {
		departureTime, err := time.ParseInLocation("2006-01-02 15:04", *req.DepartureTime, loc)
		if err != nil {
			return nil, errors.New("invalid departure_time format, use YYYY-MM-DD HH:mm")
		}
		nowInLoc := time.Now().In(loc)
		if departureTime.Before(nowInLoc) {
			return nil, errors.New("departure_time cannot be in the past")
		}
		updates["departure_time"] = departureTime
	}

	// Validasi dan set arrival_time
	if req.ArrivalTime != nil {
		arrivalTime, err := time.ParseInLocation("2006-01-02 15:04", *req.ArrivalTime, loc)
		if err != nil {
			return nil, errors.New("invalid arrival_time format, use YYYY-MM-DD HH:mm")
		}

		// Cek dengan departure_time (baik yang baru atau yang lama)
		var departureTime time.Time
		if req.DepartureTime != nil {
			departureTime, _ = time.ParseInLocation("2006-01-02 15:04", *req.DepartureTime, loc)
		} else {
			departureTime = existingSchedule.DepartureTime
		}

		if arrivalTime.Before(departureTime) || arrivalTime.Equal(departureTime) {
			return nil, errors.New("arrival_time must be after departure_time")
		}
		updates["arrival_time"] = arrivalTime
	}

	// Validasi dan set route_id
	if req.RouteID != nil {
		routeExists, err := s.scheduleRepo.CheckRouteExists(*req.RouteID)
		if err != nil {
			return nil, fmt.Errorf("error checking route: %v", err)
		}
		if !routeExists {
			return nil, errors.New("route not found")
		}
		updates["route_id"] = *req.RouteID
	}
	// Validasi dan set price
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("price must be greater than 0")
		}
		updates["price"] = *req.Price
	}

	// Update schedule
	if err := s.scheduleRepo.UpdateSchedule(id, updates); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %v", err)
	}
	// Get updated schedule untuk response
	updatedSchedule, err := s.scheduleRepo.GetScheduleByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated schedule: %v", err)
	}

	response := dto.ToScheduleResponse(*updatedSchedule, true) // true untuk admin fields
	return &response, nil
}

// DeleteSchedule - Delete schedule (Admin only)
func (s *ScheduleService) DeleteSchedule(id uint) error {
	// Cek apakah schedule exists
	existingSchedule, err := s.scheduleRepo.GetScheduleByID(id)
	if err != nil {
		return errors.New("schedule not found")
	}

	// Cek apakah ada booking yang aktif
	if existingSchedule.AvailableSeats < existingSchedule.TotalSeats {
		return errors.New("cannot delete schedule with active bookings")
	}

	// Delete schedule
	if err := s.scheduleRepo.DeleteSchedule(id); err != nil {
		return fmt.Errorf("failed to delete schedule: %v", err)
	}

	return nil
}

// GetAllSchedules - Get all schedules with pagination (Admin only)
func (s *ScheduleService) GetAllSchedules(page, limit int) (*dto.ScheduleListResponse, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 { // Batasi max limit
		limit = 100
	}
	schedules, totalCount, err := s.scheduleRepo.GetAllSchedules(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %v", err)
	}

	response := dto.ToScheduleListResponse(schedules, page, limit, totalCount, true) // true untuk admin fields
	return &response, nil
}

// SearchSchedules - Search schedules by origin, destination, and departure date (User)
func (s *ScheduleService) SearchSchedules(req dto.ScheduleSearchRequest) (*dto.ScheduleListResponse, error) {
	// Set default pagination values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	// Load timezone Indonesia (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60) // Fallback ke WIB +7
	}

	// Parse departure date dengan timezone lokal
	departureDate, err := time.ParseInLocation("2006-01-02", req.DepartureDate, loc)
	if err != nil {
		return nil, errors.New("invalid departure_date format, use YYYY-MM-DD")
	}

	// Validasi tanggal tidak boleh di masa lalu (bandingkan dengan tanggal hari ini di timezone yang sama)
	todayInLoc := time.Now().In(loc).Truncate(24 * time.Hour)
	if departureDate.Before(todayInLoc) {
		return nil, errors.New("departure_date cannot be in the past")
	}
	// Search schedules
	schedules, totalCount, err := s.scheduleRepo.SearchSchedules(req.Origin, req.Destination, departureDate, req.Page, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search schedules: %v", err)
	}

	response := dto.ToScheduleListResponse(schedules, req.Page, req.Limit, totalCount, false) // false untuk user fields
	return &response, nil
}

// GetScheduleByID - Get schedule by ID (unified untuk admin dan user)
func (s *ScheduleService) GetScheduleByID(id uint, isAdmin bool) (*dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetScheduleByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	response := dto.ToScheduleResponse(*schedule, isAdmin)
	return &response, nil
}
