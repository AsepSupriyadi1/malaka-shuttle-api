package models

// ModelHandlers untuk dependency injection dan centralized access
type ModelHandlers struct {
	User          *UserHandler
	Route         *RouteHandler
	Schedule      *ScheduleHandler
	Seat          *SeatHandler
	Booking       *BookingHandler
	BookingDetail *BookingDetailHandler
	Payment       *PaymentHandler
}

// NewModelHandlers membuat instance semua model handlers
func NewModelHandlers() *ModelHandlers {
	return &ModelHandlers{
		User:          NewUserHandler(),
		Route:         NewRouteHandler(),
		Schedule:      NewScheduleHandler(),
		Seat:          NewSeatHandler(),
		Booking:       NewBookingHandler(),
		BookingDetail: NewBookingDetailHandler(),
		Payment:       NewPaymentHandler(),
	}
}

// Global instance untuk kemudahan akses
var Handlers = NewModelHandlers()
