package cron

import (
	"log"
	"malakashuttle/services"

	"github.com/robfig/cron/v3"
)

// BookingScheduler handles scheduled tasks for booking management
type BookingScheduler struct {
	bookingService *services.BookingService
	cron           *cron.Cron
}

func NewBookingScheduler(bookingService *services.BookingService) *BookingScheduler {
	return &BookingScheduler{
		bookingService: bookingService,
		cron:           cron.New(),
	}
}

// Start starts the booking expiration scheduler using cron job
func (s *BookingScheduler) Start() {
	// Schedule to run every 10 minutes to check for expired bookings
	// Cron expression: "*/10 * * * *" means every 10 minutes
	_, err := s.cron.AddFunc("*/10 * * * *", func() {
		log.Println("Running booking expiration check...")
		if err := s.bookingService.ExpireBookings(); err != nil {
			log.Printf("Error expiring bookings: %v", err)
		} else {
			log.Println("Booking expiration check completed successfully")
		}
	})

	if err != nil {
		log.Printf("Error scheduling booking expiration job: %v", err)
		return
	}

	// Start the cron scheduler
	s.cron.Start()
	log.Println("Booking scheduler started - checking for expired bookings every 10 minutes")
}

// Stop stops the scheduler
func (s *BookingScheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("Booking scheduler stopped")
	}
}
