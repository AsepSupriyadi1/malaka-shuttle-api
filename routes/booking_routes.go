package routes

import (
	"malakashuttle/constants"
	"malakashuttle/controllers"
	"malakashuttle/middleware"

	"github.com/gin-gonic/gin"
)

func BookingRoutes(r *gin.RouterGroup, h *controllers.BookingController) {
	// User booking routes
	userRoutes := r.Group("/bookings")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_USER))
	userRoutes.GET("", h.GetUserBookings)
	userRoutes.POST("", h.CreateBooking)
	userRoutes.GET("/:id", h.GetBookingByID)
	userRoutes.GET("/:id/receipt", h.DownloadReceipt)
	userRoutes.POST("/:id/payment", h.UploadPaymentProof)

	// Admin booking routes
	adminRoutes := r.Group("/admin/bookings")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_ADMIN))
	adminRoutes.GET("", h.GetAllBookings)
	adminRoutes.GET("/:id", h.GetBookingByID)
	adminRoutes.GET("/:id/payment/download", h.DownloadPaymentProof)
	adminRoutes.PUT("/:id/status", h.UpdateBookingStatus)

	// Staff booking routes
	staffRoutes := r.Group("/staff/bookings")
	staffRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_STAFF))
	staffRoutes.GET("", h.GetAllBookings)
	staffRoutes.GET("/:id", h.GetBookingByID)
	staffRoutes.GET("/:id/payment/download", h.DownloadPaymentProof)
	staffRoutes.PUT("/:id/status", h.UpdateBookingStatus)
}
