package dto

import (
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

// BookingListResponse represents paginated booking list
type BookingListResponse struct {
	Data       []BookingResponse `json:"data"`
	Pagination PaginationMeta    `json:"pagination"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
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
