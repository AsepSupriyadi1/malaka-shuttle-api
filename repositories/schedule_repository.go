package repositories

import (
	"fmt"
	"malakashuttle/entities"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// CreateSchedule - Create new schedule with seats (menggunakan transaction)
func (r *ScheduleRepository) CreateSchedule(schedule *entities.Schedule) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create schedule
		if err := tx.Create(schedule).Error; err != nil {
			return err
		}

		// Create seats untuk schedule ini
		var seats []entities.Seat
		for i := 1; i <= schedule.TotalSeats; i++ {
			seat := entities.Seat{
				ScheduleID: schedule.ID,
				SeatNumber: generateSeatNumber(i),
				IsBooked:   false,
			}
			seats = append(seats, seat)
		}

		// Batch insert seats
		if err := tx.CreateInBatches(seats, 100).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetScheduleByID - Get schedule by ID with route relation
func (r *ScheduleRepository) GetScheduleByID(id uint) (*entities.Schedule, error) {
	var schedule entities.Schedule
	err := r.db.Preload("Route").First(&schedule, id).Error
	if err != nil {
		return nil, err
	}

	// Hitung available_seats secara real-time berdasarkan seat yang is_booked = false
	if err := r.updateAvailableSeats(&schedule); err != nil {
		return nil, err
	}

	return &schedule, nil
}

// updateAvailableSeats - Update available_seats berdasarkan seat yang is_booked = false
func (r *ScheduleRepository) updateAvailableSeats(schedule *entities.Schedule) error {
	var availableCount int64

	// Hitung seat yang is_booked = false
	if err := r.db.Model(&entities.Seat{}).
		Where("schedule_id = ? AND is_booked = ?", schedule.ID, false).
		Count(&availableCount).Error; err != nil {
		return err
	}

	// Update available_seats di database dan struct
	schedule.AvailableSeats = int(availableCount)
	if err := r.db.Model(schedule).Update("available_seats", availableCount).Error; err != nil {
		return err
	}

	return nil
}

// UpdateSchedule - Update schedule
// Note: TotalSeats tidak bisa diubah untuk menghindari konflik data
func (r *ScheduleRepository) UpdateSchedule(id uint, updates map[string]interface{}) error {
	// Pastikan total_seats tidak ada dalam updates
	delete(updates, "total_seats")
	delete(updates, "available_seats")

	// Update schedule langsung tanpa transaksi kompleks
	if err := r.db.Model(&entities.Schedule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DeleteSchedule - Delete schedule (soft delete)
func (r *ScheduleRepository) DeleteSchedule(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Soft delete seats
		if err := tx.Where("schedule_id = ?", id).Delete(&entities.Seat{}).Error; err != nil {
			return err
		}

		// Soft delete schedule
		if err := tx.Delete(&entities.Schedule{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetAllSchedules - Get all schedules with pagination (Admin)
func (r *ScheduleRepository) GetAllSchedules(page, limit int) ([]entities.Schedule, int64, error) {
	var schedules []entities.Schedule
	var totalCount int64

	// Count total records
	if err := r.db.Model(&entities.Schedule{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	// Get paginated results
	offset := (page - 1) * limit
	err := r.db.Preload("Route").
		Order("departure_time ASC").
		Limit(limit).
		Offset(offset).
		Find(&schedules).Error

	if err != nil {
		return nil, 0, err
	}

	// Update available_seats untuk setiap schedule
	for i := range schedules {
		if err := r.updateAvailableSeats(&schedules[i]); err != nil {
			return nil, 0, err
		}
	}

	return schedules, totalCount, nil
}

// SearchSchedules - Search schedules by origin, destination, and departure date (User)
func (r *ScheduleRepository) SearchSchedules(origin, destination string, departureDate time.Time, page, limit int) ([]entities.Schedule, int64, error) {
	var schedules []entities.Schedule
	var totalCount int64

	// Parse date range untuk mencari schedule pada hari tertentu
	startOfDay := time.Date(departureDate.Year(), departureDate.Month(), departureDate.Day(), 0, 0, 0, 0, departureDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query builder - perbaikan nama kolom
	query := r.db.Model(&entities.Schedule{}).
		Joins("JOIN routes ON routes.id = schedules.route_id").
		Where("LOWER(routes.origin_city) = LOWER(?)", origin).
		Where("LOWER(routes.destination_city) = LOWER(?)", destination).
		Where("schedules.departure_time >= ? AND schedules.departure_time < ?", startOfDay, endOfDay).
		Where("schedules.available_seats > 0") // Hanya tampilkan yang masih ada seat tersedia

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	// Get paginated results
	offset := (page - 1) * limit
	err := query.Preload("Route").
		Order("schedules.departure_time ASC"). // Urutkan dari paling pagi
		Limit(limit).
		Offset(offset).
		Find(&schedules).Error
	if err != nil {
		return nil, 0, err
	}

	// Update available_seats untuk setiap schedule
	for i := range schedules {
		if err := r.updateAvailableSeats(&schedules[i]); err != nil {
			return nil, 0, err
		}
	}

	return schedules, totalCount, nil
}

// CheckRouteExists - Check if route exists
func (r *ScheduleRepository) CheckRouteExists(routeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entities.Route{}).Where("id = ?", routeID).Count(&count).Error
	return count > 0, err
}

// GetSeatsByScheduleID - Get all seats for a schedule
func (r *ScheduleRepository) GetSeatsByScheduleID(scheduleID uint) ([]entities.Seat, error) {
	var seats []entities.Seat
	err := r.db.Where("schedule_id = ?", scheduleID).Order("seat_number").Find(&seats).Error
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// generateSeatNumber - Generate seat number (A1, A2, B1, B2, ...)
func generateSeatNumber(seatIndex int) string {
	// Buat format seat number seperti A1, A2, A3, A4, B1, B2, dst
	// Asumsi 4 seat per row
	row := ((seatIndex - 1) / 4)
	col := ((seatIndex - 1) % 4) + 1

	rowLetter := string(rune('A' + row))
	return fmt.Sprintf("%s%d", rowLetter, col)
}
