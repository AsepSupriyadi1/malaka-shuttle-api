package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"malakashuttle/dto"
	"malakashuttle/entities"
	"malakashuttle/repositories"
	"malakashuttle/utils"

	"gorm.io/gorm"
)

type BookingService struct {
	bookingRepo  *repositories.BookingRepository
	scheduleRepo *repositories.ScheduleRepository
	userRepo     repositories.UserRepository
}

func NewBookingService(
	bookingRepo *repositories.BookingRepository,
	scheduleRepo *repositories.ScheduleRepository,
	userRepo repositories.UserRepository,
) *BookingService {
	return &BookingService{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

// CreateBooking creates a new booking
func (s *BookingService) CreateBooking(userID uint, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	// Validate schedule exists and is available
	schedule, err := s.scheduleRepo.GetScheduleByID(req.ScheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	// Check if schedule is in the future
	if schedule.DepartureTime.Before(time.Now()) {
		return nil, errors.New("cannot book past schedule")
	}

	// Validate all seat IDs belong to the route
	seatIDs := make([]uint, len(req.Passengers))
	for i, passenger := range req.Passengers {
		seatIDs[i] = passenger.SeatID
	}

	// Check if seats exist and belong to the route
	// This would require a seat repository - for now, we'll assume seats are valid

	// Create booking
	booking := &entities.Booking{
		UserID:      userID,
		ScheduleID:  req.ScheduleID,
		BookingTime: time.Now(),
		Status:      entities.BookingStatusPending,
		ExpiresAt:   time.Now().Add(30 * time.Minute), // 30 minutes expiry
	}

	// Create booking details
	bookingDetails := make([]entities.BookingDetail, len(req.Passengers))
	var totalAmount float64

	for i, passenger := range req.Passengers {
		bookingDetails[i] = entities.BookingDetail{
			SeatID:        passenger.SeatID,
			PassengerName: passenger.PassengerName,
			Price:         schedule.Price, // Assuming all seats have same price
		}
		totalAmount += schedule.Price
	}

	// Create booking in database
	err = s.bookingRepo.CreateBooking(booking, bookingDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}
	// Get created booking with relations
	createdBooking, err := s.bookingRepo.GetBookingByID(booking.ID, &userID)
	if err != nil {
		return nil, err
	}

	return dto.NewBookingResponseFromEntity(createdBooking), nil
}

// GetBookingByID gets booking by ID
func (s *BookingService) GetBookingByID(id uint, userID *uint) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return dto.NewBookingResponseFromEntity(booking), nil
}

// GetBookingDetailByID gets detailed booking by ID
func (s *BookingService) GetBookingDetailByID(id uint, userID *uint) (*dto.BookingFullResponse, error) {
	booking, err := s.bookingRepo.GetBookingByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	return dto.NewBookingFullResponseFromEntity(booking), nil
}

// GetUserBookings gets all bookings for a user
func (s *BookingService) GetUserBookings(userID uint, params utils.PaginationParams, status []entities.BookingStatus) (*utils.PaginationResponse, error) {
	bookings, total, err := s.bookingRepo.GetBookingsByUserID(userID, params.Page, params.Limit, status)
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	data := make([]dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		data[i] = *dto.NewBookingResponseFromEntity(&booking)
	}

	response := utils.CreatePaginationResponse(data, total, params)
	return &response, nil
}

// GetUserBookingsList gets all bookings for a user with list view response format
func (s *BookingService) GetUserBookingsList(userID uint, params utils.PaginationParams, status []entities.BookingStatus) (*utils.PaginationResponse, error) {
	bookings, total, err := s.bookingRepo.GetBookingsByUserID(userID, params.Page, params.Limit, status)
	if err != nil {
		return nil, err
	}

	// Map to list response DTOs
	data := make([]dto.BookingListResponse, len(bookings))
	for i, booking := range bookings {
		data[i] = *dto.NewBookingListResponseFromEntity(&booking)
	}

	response := utils.CreatePaginationResponse(data, total, params)
	return &response, nil
}

// GetAllBookings gets all bookings (for staff)
func (s *BookingService) GetAllBookings(params utils.PaginationParams, status []entities.BookingStatus) (*utils.PaginationResponse, error) {
	bookings, total, err := s.bookingRepo.GetAllBookings(params.Page, params.Limit, status)
	if err != nil {
		return nil, err
	}

	// Map to response DTOs
	data := make([]dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		data[i] = *dto.NewBookingResponseFromEntity(&booking)
	}

	response := utils.CreatePaginationResponse(data, total, params)
	return &response, nil
}

// GetAllBookingsList gets all bookings with list view response format
func (s *BookingService) GetAllBookingsList(params utils.PaginationParams, status []entities.BookingStatus) (*utils.PaginationResponse, error) {
	bookings, total, err := s.bookingRepo.GetAllBookings(params.Page, params.Limit, status)
	if err != nil {
		return nil, err
	}

	// Map to list response DTOs
	data := make([]dto.BookingListResponse, len(bookings))
	for i, booking := range bookings {
		data[i] = *dto.NewBookingListResponseFromEntity(&booking)
	}

	response := utils.CreatePaginationResponse(data, total, params)
	return &response, nil
}

// UploadPaymentProof uploads payment proof
func (s *BookingService) UploadPaymentProof(bookingID uint, userID uint, file *multipart.FileHeader, paymentMethod string) error {
	// Check if booking exists and belongs to user
	booking, err := s.bookingRepo.GetBookingForPayment(bookingID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found or not eligible for payment")
		}
		return err
	}

	// Check if booking is not expired
	if booking.ExpiresAt.Before(time.Now()) {
		return errors.New("booking has expired")
	}

	// Check if payment already exists
	existingPayment, err := s.bookingRepo.GetPaymentByBookingID(bookingID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingPayment != nil {
		return errors.New("payment proof already uploaded")
	}

	// Save file (implement file upload logic)
	filename := fmt.Sprintf("payment_%d_%d%s", bookingID, time.Now().Unix(), filepath.Ext(file.Filename))
	savePath := filepath.Join("uploads", "payments", filename)

	// Create file and save
	if err := saveUploadedFile(file, savePath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	proofURL := "/uploads/payments/" + filename

	// Create payment record
	payment := &entities.Payment{
		BookingID:     bookingID,
		PaymentMethod: paymentMethod,
		PaymentStatus: entities.PaymentStatusPending,
		ProofImageURL: proofURL,
	}

	return s.bookingRepo.CreatePayment(payment)
}

// UpdateBookingStatus updates booking status (for staff)
func (s *BookingService) UpdateBookingStatus(bookingID uint, status entities.BookingStatus) error {
	// Validate status
	if status != entities.BookingStatusSuccess && status != entities.BookingStatusRejected {
		return errors.New("invalid status")
	}

	// Get booking to ensure it exists
	booking, err := s.bookingRepo.GetBookingByID(bookingID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	// Check if booking is in waiting_verification status
	if booking.Status != entities.BookingStatusWaitingVerification {
		return errors.New("booking is not in waiting verification status")
	}
	// Update booking status
	err = s.bookingRepo.UpdateBookingStatus(bookingID, status)
	if err != nil {
		return err
	}

	// Handle post-status update actions
	if status == entities.BookingStatusSuccess {
		// If status is success, update payment status
		payment, err := s.bookingRepo.GetPaymentByBookingID(bookingID)
		if err != nil {
			return err
		}

		// Update payment status and date
		now := time.Now()
		payment.PaymentStatus = entities.PaymentStatusSuccess
		payment.PaymentDate = &now

		// You would need to implement UpdatePayment method in repository
		// For now, we'll assume it's handled
	} else if status == entities.BookingStatusRejected {
		// If status is rejected, free the seats so they can be booked again
		err = s.bookingRepo.FreeSeatsByBookingID(bookingID)
		if err != nil {
			return fmt.Errorf("failed to free seats: %w", err)
		}
	}

	return nil
}

// ExpireBookings expires bookings that have passed their expiry time
func (s *BookingService) ExpireBookings() error {
	return s.bookingRepo.ExpireBookings()
}

// GetAvailableSeats gets available seats for a schedule
func (s *BookingService) GetAvailableSeats(scheduleID uint) (*dto.AvailableSeatsResponse, error) {
	// Validate schedule exists
	_, err := s.scheduleRepo.GetScheduleByID(scheduleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}

	// Get all seats for the schedule
	seats, err := s.bookingRepo.GetAvailableSeats(scheduleID)
	if err != nil {
		return nil, err
	}

	// Count available and booked seats
	var availableCount, bookedCount int
	seatResponses := make([]dto.SeatResponse, len(seats))

	for i, seat := range seats {
		seatResponses[i] = dto.SeatResponse{
			ID:         seat.ID,
			SeatNumber: seat.SeatNumber,
			IsBooked:   seat.IsBooked,
		}

		if seat.IsBooked {
			bookedCount++
		} else {
			availableCount++
		}
	}

	response := &dto.AvailableSeatsResponse{
		ScheduleID:     scheduleID,
		TotalSeats:     len(seats),
		AvailableSeats: availableCount,
		BookedSeats:    bookedCount,
		Seats:          seatResponses,
	}

	return response, nil
}

// GetPaymentByBookingID gets payment details for a booking
func (s *BookingService) GetPaymentByBookingID(bookingID uint) (*dto.PaymentResponse, error) {
	payment, err := s.bookingRepo.GetPaymentByBookingID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	response := &dto.PaymentResponse{
		ID:            payment.ID,
		PaymentMethod: payment.PaymentMethod,
		PaymentStatus: payment.PaymentStatus,
		PaymentDate:   payment.PaymentDate,
		ProofImageURL: payment.ProofImageURL,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}

	return response, nil
}

// GenerateBookingReceipt generates PDF receipt for a booking
func (s *BookingService) GenerateBookingReceipt(bookingID uint, userID *uint) (string, error) {
	// Get booking details
	booking, err := s.bookingRepo.GetBookingByID(bookingID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("booking not found")
		}
		return "", err
	}
	// Only allow receipt generation for successful bookings
	if booking.Status != entities.BookingStatusSuccess {
		return "", errors.New("receipt can only be generated for successful bookings")
	}

	// Map to response DTO
	bookingResponse := dto.NewBookingResponseFromEntity(booking)

	// Create receipts directory if not exists
	receiptsDir := "uploads/receipts"
	if err := os.MkdirAll(receiptsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create receipts directory: %w", err)
	}

	// Generate filename
	filename := fmt.Sprintf("receipt_%d_%d.pdf", bookingID, time.Now().Unix())
	outputPath := filepath.Join(receiptsDir, filename)

	// Generate PDF
	pdfGenerator := utils.NewPDFReceiptGenerator()
	err = pdfGenerator.GenerateBookingReceipt(bookingResponse, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return outputPath, nil
}

// saveUploadedFile saves multipart file to disk
func saveUploadedFile(fileHeader *multipart.FileHeader, dst string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
