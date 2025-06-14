package routes

import (
	"malakashuttle/controllers"
	"malakashuttle/middleware"
	"malakashuttle/repositories"
	"malakashuttle/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	routeRepo := repositories.NewRouteRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	bookingRepo := repositories.NewBookingRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	routeService := services.NewRouteService(routeRepo)
	scheduleService := services.NewScheduleService(scheduleRepo)
	bookingService := services.NewBookingService(bookingRepo, scheduleRepo, userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	routeController := controllers.NewRouteController(routeService)
	scheduleController := controllers.NewScheduleController(scheduleService)
	bookingController := controllers.NewBookingController(bookingService)
	testController := controllers.NewTestController()

	// Apply logging middleware to all API routes
	api := router.Group("/api")
	api.Use(middleware.LoggerMiddleware()) // Apply logging middleware at API level
	{
		// Test routes (no authentication needed)
		api.GET("/ping", testController.Ping)

		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
		// Public routes (no authentication needed)
		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware()) // Apply JWT middleware
		protected.Use(middleware.RequireUser(db))
		{
			// Schedule routes for public access
			protected.GET("/schedules/search", scheduleController.SearchSchedules)     // GET /api/public/schedules/search
			protected.GET("/schedules/:id", scheduleController.GetScheduleByIDForUser) // GET /api/public/schedules/:id
			protected.GET("/schedules/:id/seats", bookingController.GetAvailableSeats) // GET /api/protected/schedules/:id/seats

			// Booking routes for authenticated users
			bookings := protected.Group("/bookings")
			{
				bookings.POST("", bookingController.CreateBooking)                  // POST /api/protected/bookings
				bookings.GET("", bookingController.GetUserBookings)                 // GET /api/protected/bookings
				bookings.GET("/:id", bookingController.GetBookingByID)              // GET /api/protected/bookings/:id
				bookings.POST("/:id/payment", bookingController.UploadPaymentProof) // POST /api/protected/bookings/:id/payment
			}
		}

		// Admin routes (require admin role)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware()) // Apply JWT middleware
		admin.Use(middleware.RequireAdmin(db)) // Apply admin role middleware
		{
			adminRoutes := admin.Group("/routes")
			{
				adminRoutes.GET("", routeController.GetAllRoutes)       // GET /api/admin/routes
				adminRoutes.POST("", routeController.CreateRoute)       // POST /api/admin/routes				adminRoutes.GET("/:id", routeController.GetRouteByID) // GET /api/admin/routes/:id
				adminRoutes.PUT("/:id", routeController.UpdateRoute)    // PUT /api/admin/routes/:id
				adminRoutes.DELETE("/:id", routeController.DeleteRoute) // DELETE /api/admin/routes/:id
			}

			adminSchedules := admin.Group("/schedules")
			{
				adminSchedules.GET("", scheduleController.GetAllSchedules)       // GET /api/admin/schedules
				adminSchedules.POST("", scheduleController.CreateSchedule)       // POST /api/admin/schedules
				adminSchedules.GET("/:id", scheduleController.GetScheduleByID)   // GET /api/admin/schedules/:id (untuk admin)
				adminSchedules.PUT("/:id", scheduleController.UpdateSchedule)    // PUT /api/admin/schedules/:id
				adminSchedules.DELETE("/:id", scheduleController.DeleteSchedule) // DELETE /api/admin/schedules/:id
			} // Admin booking routes
			adminBookings := admin.Group("/bookings")
			{
				adminBookings.GET("", bookingController.GetAllBookings)                            // GET /api/admin/bookings
				adminBookings.GET("/:id", bookingController.GetBookingByID)                        // GET /api/admin/bookings/:id
				adminBookings.GET("/:id/payment/download", bookingController.DownloadPaymentProof) // GET /api/admin/bookings/:id/payment/download
				adminBookings.PUT("/:id/status", bookingController.UpdateBookingStatus)            // PUT /api/admin/bookings/:id/status
			}
		}
	}
}
