package repositories

import (
	"errors"
	"time"

	"malakashuttle/entities"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// CreateBooking creates a new booking with booking details in a transaction
func (r *BookingRepository) CreateBooking(booking *entities.Booking, bookingDetails []entities.BookingDetail) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if any seats are already booked for this schedule
		var existingDetails []entities.BookingDetail
		seatIDs := make([]uint, len(bookingDetails))
		for i, detail := range bookingDetails {
			seatIDs[i] = detail.SeatID
		}

		// Check for active booking details (not deleted) that conflict
		err := tx.Joins("JOIN bookings ON booking_details.booking_id = bookings.id").
			Where("booking_details.seat_id IN ? AND bookings.schedule_id = ? AND bookings.status NOT IN ? AND booking_details.deleted_at IS NULL",
				seatIDs, booking.ScheduleID, []entities.BookingStatus{
					entities.BookingStatusExpired,
					entities.BookingStatusCancelled,
					entities.BookingStatusRejected,
				}).
			Find(&existingDetails).Error

		if err != nil {
			return err
		}

		if len(existingDetails) > 0 {
			return errors.New("one or more seats are already booked")
		}

		// Double check seat availability
		var bookedSeats []entities.Seat
		err = tx.Where("id IN ? AND is_booked = ?", seatIDs, true).Find(&bookedSeats).Error
		if err != nil {
			return err
		}

		if len(bookedSeats) > 0 {
			return errors.New("one or more seats are already booked")
		}

		// Create booking
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		// Set booking ID for details and create them
		for i := range bookingDetails {
			bookingDetails[i].BookingID = booking.ID
		}

		if err := tx.Create(&bookingDetails).Error; err != nil {
			return err
		}

		// Update seat status to booked
		if err := tx.Model(&entities.Seat{}).Where("id IN ?", seatIDs).Update("is_booked", true).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetBookingByID retrieves a booking by ID with all relations
func (r *BookingRepository) GetBookingByID(id uint, userID *uint) (*entities.Booking, error) {
	var booking entities.Booking
	query := r.db.Preload("Schedule").
		Preload("Schedule.Route").
		Preload("BookingDetails").
		Preload("BookingDetails.Seat").
		Preload("Payment").
		Preload("User")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.First(&booking, id).Error
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// GetBookingsByUserID retrieves all bookings for a specific user
func (r *BookingRepository) GetBookingsByUserID(userID uint, page, limit int, status []entities.BookingStatus) ([]entities.Booking, int64, error) {
	var bookings []entities.Booking
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	// Count total
	if err := query.Model(&entities.Booking{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * limit
	err := query.Preload("Schedule").
		Preload("Schedule.Route").
		Preload("BookingDetails").
		Preload("Payment").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error

	return bookings, total, err
}

// GetAllBookings retrieves all bookings (for staff)
func (r *BookingRepository) GetAllBookings(page, limit int, status []entities.BookingStatus) ([]entities.Booking, int64, error) {
	var bookings []entities.Booking
	var total int64

	query := r.db.Model(&entities.Booking{})

	if len(status) > 0 {
		query = query.Where("status IN ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// Get paginated data
	offset := (page - 1) * limit
	err := query.Preload("Schedule").
		Preload("Schedule.Route").
		Preload("BookingDetails").
		Preload("Payment").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error

	return bookings, total, err
}

// UpdateBookingStatus updates booking status
func (r *BookingRepository) UpdateBookingStatus(id uint, status entities.BookingStatus) error {
	return r.db.Model(&entities.Booking{}).Where("id = ?", id).Update("status", status).Error
}

// ExpireBookings updates status of bookings that have expired
func (r *BookingRepository) ExpireBookings() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get expired bookings with their seat IDs
		var expiredBookings []entities.Booking
		err := tx.Preload("BookingDetails").
			Where("expires_at <= ? AND status = ?", time.Now(), entities.BookingStatusPending).
			Find(&expiredBookings).Error
		if err != nil {
			return err
		}

		if len(expiredBookings) == 0 {
			return nil // No expired bookings
		}

		// Collect all seat IDs from expired bookings
		var seatIDs []uint
		var bookingIDs []uint
		for _, booking := range expiredBookings {
			bookingIDs = append(bookingIDs, booking.ID)
			for _, detail := range booking.BookingDetails {
				seatIDs = append(seatIDs, detail.SeatID)
			}
		}
		// Update booking status to expired
		if err := tx.Model(&entities.Booking{}).Where("id IN ?", bookingIDs).
			Update("status", entities.BookingStatusExpired).Error; err != nil {
			return err
		}

		// Soft delete booking details to free up unique constraint
		if err := tx.Where("booking_id IN ?", bookingIDs).Delete(&entities.BookingDetail{}).Error; err != nil {
			return err
		}

		// Free up the seats
		if len(seatIDs) > 0 {
			if err := tx.Model(&entities.Seat{}).Where("id IN ?", seatIDs).
				Update("is_booked", false).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetBookingForPayment gets booking for payment (only pending bookings)
func (r *BookingRepository) GetBookingForPayment(id uint, userID uint) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("id = ? AND user_id = ? AND status = ?", id, userID, entities.BookingStatusPending).
		First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// CreatePayment creates payment record
func (r *BookingRepository) CreatePayment(payment *entities.Payment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create payment
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// Update booking status to waiting_verification
		if err := tx.Model(&entities.Booking{}).Where("id = ?", payment.BookingID).
			Update("status", entities.BookingStatusWaitingVerification).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetPaymentByBookingID gets payment by booking ID
func (r *BookingRepository) GetPaymentByBookingID(bookingID uint) (*entities.Payment, error) {
	var payment entities.Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetAvailableSeats gets all seats for a schedule with their booking status
func (r *BookingRepository) GetAvailableSeats(scheduleID uint) ([]entities.Seat, error) {
	var seats []entities.Seat
	err := r.db.Where("schedule_id = ?", scheduleID).Order("seat_number").Find(&seats).Error
	return seats, err
}

// FreeSeatsByBookingID frees seats when booking is rejected/cancelled
func (r *BookingRepository) FreeSeatsByBookingID(bookingID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get seat IDs from booking details
		var bookingDetails []entities.BookingDetail
		err := tx.Where("booking_id = ?", bookingID).Find(&bookingDetails).Error
		if err != nil {
			return err
		}

		// Extract seat IDs
		var seatIDs []uint
		for _, detail := range bookingDetails {
			seatIDs = append(seatIDs, detail.SeatID)
		}

		// Soft delete booking details (this will free up the unique constraint)
		if err := tx.Where("booking_id = ?", bookingID).Delete(&entities.BookingDetail{}).Error; err != nil {
			return err
		}

		// Free the seats
		if len(seatIDs) > 0 {
			err = tx.Model(&entities.Seat{}).Where("id IN ?", seatIDs).Update("is_booked", false).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}
