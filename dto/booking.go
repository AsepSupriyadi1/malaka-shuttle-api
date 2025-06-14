package dto

import (
	"fmt"
	"time"

	"malakashuttle/entities"
)

// BookingPassenger represents passenger data for booking
type BookingPassenger struct {
	PassengerName string `json:"passenger_name" validate:"required,min=2,max=100"`
	SeatID        uint   `json:"seat_id" validate:"required,min=1"`
}

// CreateBookingRequest represents the request payload for creating a booking
type CreateBookingRequest struct {
	ScheduleID uint               `json:"schedule_id" validate:"required,min=1"`
	Passengers []BookingPassenger `json:"passengers" validate:"required,min=1,max=10,dive"`
}

// BookingResponse represents booking data in response
type BookingResponse struct {
	ID          uint                   `json:"id"`
	Status      entities.BookingStatus `json:"status"`
	ExpiresAt   time.Time              `json:"expires_at"`
	TotalAmount float64                `json:"total_amount"`
	Schedule    *ScheduleResponse      `json:"schedule,omitempty"`
	Passengers  []PassengerResponse    `json:"passengers,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// PassengerResponse represents passenger data in response
type PassengerResponse struct {
	PassengerName string `json:"passenger_name"`
	SeatNumber    string `json:"seat_number"`
}

// BookingDetailResponse represents booking detail in response
type BookingDetailResponse struct {
	ID            uint    `json:"id"`
	SeatID        uint    `json:"seat_id"`
	PassengerName string  `json:"passenger_name"`
	Price         float64 `json:"price"`
	SeatNumber    string  `json:"seat_number,omitempty"`
}

// PaymentResponse represents payment data in response
type PaymentResponse struct {
	ID            uint                   `json:"id"`
	PaymentMethod string                 `json:"payment_method"`
	PaymentStatus entities.PaymentStatus `json:"payment_status"`
	PaymentDate   *time.Time             `json:"payment_date"`
	ProofImageURL string                 `json:"proof_image_url,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// UploadPaymentProofRequest represents the request for uploading payment proof
type UploadPaymentProofRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,min=2,max=50"`
}

// UpdateBookingStatusRequest represents the request for updating booking status (staff only)
type UpdateBookingStatusRequest struct {
	Status entities.BookingStatus `json:"status" validate:"required,oneof=success rejected"`
	Notes  string                 `json:"notes,omitempty" validate:"max=500"`
}

// SeatResponse represents seat data for available seats endpoint
type SeatResponse struct {
	ID         uint   `json:"id"`
	SeatNumber string `json:"seat_number"`
	IsBooked   bool   `json:"is_booked"`
}

// AvailableSeatsResponse represents available seats for a schedule
type AvailableSeatsResponse struct {
	ScheduleID     uint           `json:"schedule_id"`
	TotalSeats     int            `json:"total_seats"`
	AvailableSeats int            `json:"available_seats"`
	BookedSeats    int            `json:"booked_seats"`
	Seats          []SeatResponse `json:"seats"`
}

// BookingListResponse represents a booking in list view format
type BookingListResponse struct {
	BookingID      uint                   `json:"booking_id"`
	BookingStatus  entities.BookingStatus `json:"booking_status"`
	ExpiresAt      string                 `json:"expires_at"`
	TotalAmount    float64                `json:"total_amount"`
	Origin         string                 `json:"origin"`
	Destination    string                 `json:"destination"`
	DepartureTime  string                 `json:"departure_time"`
	ArrivalTime    string                 `json:"arrival_time"`
	Duration       string                 `json:"duration,omitempty"`
	PassengerCount int                    `json:"passenger_count,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

// BookingFullResponse represents detailed booking data for single booking view
type BookingFullResponse struct {
	BookingID        uint                      `json:"booking_id"`
	BookingStatus    entities.BookingStatus    `json:"booking_status"`
	ExpiresAt        string                    `json:"expires_at"`
	TotalAmount      float64                   `json:"total_amount"`
	Origin           string                    `json:"origin"`
	Destination      string                    `json:"destination"`
	DepartureTime    string                    `json:"departure_time"`
	ArrivalTime      string                    `json:"arrival_time"`
	Duration         string                    `json:"duration"`
	Price            float64                   `json:"price"`
	PassengerDetails []PassengerDetailResponse `json:"passenger_details"`
	PaymentInfo      *PaymentInfoResponse      `json:"payment_info,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
}

// PassengerDetailResponse represents detailed passenger data
type PassengerDetailResponse struct {
	PassengerName string  `json:"passenger_name"`
	SeatNumber    string  `json:"seat_number"`
	Price         float64 `json:"price"`
}

// PaymentInfoResponse represents payment information
type PaymentInfoResponse struct {
	PaymentMethod string                 `json:"payment_method,omitempty"`
	PaymentStatus entities.PaymentStatus `json:"payment_status"`
	PaymentDate   *time.Time             `json:"payment_date,omitempty"`
	ProofImageURL string                 `json:"proof_image_url,omitempty"`
	AdminNotes    string                 `json:"admin_notes,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// Helper function to calculate duration between two times
func calculateDuration(departure, arrival time.Time) string {
	duration := arrival.Sub(departure)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// FromEntity creates a BookingResponse from a Booking entity
func (b *BookingResponse) FromEntity(booking *entities.Booking) {
	b.ID = booking.ID
	b.Status = booking.Status
	b.ExpiresAt = booking.ExpiresAt
	b.CreatedAt = booking.CreatedAt
	// Map schedule if loaded
	if booking.Schedule.ID != 0 {
		b.Schedule = &ScheduleResponse{
			ID:            booking.Schedule.ID,
			Origin:        booking.Schedule.Route.OriginCity,
			Destination:   booking.Schedule.Route.DestinationCity,
			DepartureTime: booking.Schedule.DepartureTime.Format("2006-01-02 15:04"),
			ArrivalTime:   booking.Schedule.ArrivalTime.Format("2006-01-02 15:04"),
			Price:         booking.Schedule.Price,
			Duration:      calculateDuration(booking.Schedule.DepartureTime, booking.Schedule.ArrivalTime),
			CreatedAt:     &booking.Schedule.CreatedAt,
			UpdatedAt:     &booking.Schedule.UpdatedAt,
		}
	}
	// Map passengers (simplified from booking details)
	b.Passengers = make([]PassengerResponse, len(booking.BookingDetails))
	for i, detail := range booking.BookingDetails {
		b.Passengers[i] = PassengerResponse{
			PassengerName: detail.PassengerName,
		}

		// Add seat number if seat is loaded
		if detail.Seat.ID != 0 {
			b.Passengers[i].SeatNumber = detail.Seat.SeatNumber
		}
	}

	// Use the pre-calculated payment amount from the booking entity
	b.TotalAmount = booking.PaymentAmount
}

// NewBookingResponseFromEntity creates a new BookingResponse from a Booking entity
func NewBookingResponseFromEntity(booking *entities.Booking) *BookingResponse {
	response := &BookingResponse{}
	response.FromEntity(booking)
	return response
}

// FromEntity creates a BookingListResponse from a Booking entity
func (b *BookingListResponse) FromEntity(booking *entities.Booking) {
	b.BookingID = booking.ID
	b.BookingStatus = booking.Status
	b.ExpiresAt = booking.ExpiresAt.Format("2006-01-02 15:04")
	b.CreatedAt = booking.CreatedAt
	// Map schedule data directly into response fields
	if booking.Schedule.ID != 0 {
		b.Origin = booking.Schedule.Route.OriginCity
		b.Destination = booking.Schedule.Route.DestinationCity
		b.DepartureTime = booking.Schedule.DepartureTime.Format("2006-01-02 15:04")
		b.ArrivalTime = booking.Schedule.ArrivalTime.Format("2006-01-02 15:04")
		b.Duration = calculateDuration(booking.Schedule.DepartureTime, booking.Schedule.ArrivalTime)
	}
	// Set passenger count and use pre-calculated payment amount
	b.PassengerCount = len(booking.BookingDetails)
	b.TotalAmount = booking.PaymentAmount
}

// NewBookingListResponseFromEntity creates a new BookingListResponse from a Booking entity
func NewBookingListResponseFromEntity(booking *entities.Booking) *BookingListResponse {
	response := &BookingListResponse{}
	response.FromEntity(booking)
	return response
}

// FromEntity creates a BookingFullResponse from a Booking entity
func (b *BookingFullResponse) FromEntity(booking *entities.Booking) {
	b.BookingID = booking.ID
	b.BookingStatus = booking.Status
	b.ExpiresAt = booking.ExpiresAt.Format("2006-01-02 15:04")
	b.CreatedAt = booking.CreatedAt
	b.UpdatedAt = booking.UpdatedAt
	// Map schedule data directly into response fields
	if booking.Schedule.ID != 0 {
		b.Origin = booking.Schedule.Route.OriginCity
		b.Destination = booking.Schedule.Route.DestinationCity
		b.DepartureTime = booking.Schedule.DepartureTime.Format("2006-01-02 15:04")
		b.ArrivalTime = booking.Schedule.ArrivalTime.Format("2006-01-02 15:04")
		b.Price = booking.Schedule.Price
		b.Duration = calculateDuration(booking.Schedule.DepartureTime, booking.Schedule.ArrivalTime)
	}
	// Map passenger details with seat information
	b.PassengerDetails = make([]PassengerDetailResponse, len(booking.BookingDetails))
	for i, detail := range booking.BookingDetails {
		b.PassengerDetails[i] = PassengerDetailResponse{
			PassengerName: detail.PassengerName,
			Price:         detail.Price,
		}

		// Add seat number if seat is loaded
		if detail.Seat.ID != 0 {
			b.PassengerDetails[i].SeatNumber = detail.Seat.SeatNumber
		}
	}

	// Use pre-calculated payment amount from booking entity
	b.TotalAmount = booking.PaymentAmount

	// Map payment information if available
	if booking.Payment != nil {
		b.PaymentInfo = &PaymentInfoResponse{
			PaymentMethod: booking.Payment.PaymentMethod,
			PaymentStatus: booking.Payment.PaymentStatus,
			PaymentDate:   booking.Payment.PaymentDate,
			ProofImageURL: booking.Payment.ProofImageURL,
			CreatedAt:     booking.Payment.CreatedAt,
			UpdatedAt:     booking.Payment.UpdatedAt,
		}
	}
}

// NewBookingFullResponseFromEntity creates a new BookingFullResponse from a Booking entity
func NewBookingFullResponseFromEntity(booking *entities.Booking) *BookingFullResponse {
	response := &BookingFullResponse{}
	response.FromEntity(booking)
	return response
}
