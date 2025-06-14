package services

import (
	"log"
	"time"
)

// BookingScheduler handles scheduled tasks for booking management
type BookingScheduler struct {
	bookingService *BookingService
}

func NewBookingScheduler(bookingService *BookingService) *BookingScheduler {
	return &BookingScheduler{
		bookingService: bookingService,
	}
}

// StartScheduler starts the booking expiration scheduler
func (s *BookingScheduler) Start() {
	// Run every minute to check for expired bookings
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			if err := s.bookingService.ExpireBookings(); err != nil {
				log.Printf("Error expiring bookings: %v", err)
			}
		}
	}()

	log.Println("Booking scheduler started - checking for expired bookings every minute")
}

// Stop stops the scheduler (you can implement this if needed)
func (s *BookingScheduler) Stop() {
	// Implementation to stop the scheduler
}
