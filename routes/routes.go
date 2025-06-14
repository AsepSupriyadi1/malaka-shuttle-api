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

		protected := api.Group("/protected")
		protected.Use(middleware.AuthMiddleware()) // Apply JWT middleware
		{
			// User routes (require authentication)
			user := protected.Group("/user")
			user.Use(middleware.RequireUser(db)) // Apply user role middleware
			{
				// Schedule routes for user access
				user.GET("/schedules/search", scheduleController.SearchSchedules)
				user.GET("/schedule/:id", scheduleController.GetScheduleByID)

				// Booking routes for authenticated users
				bookings := user.Group("/bookings")
				{
					bookings.GET("", bookingController.GetUserBookings)
				}

				booking := user.Group("/booking")
				{
					booking.POST("", bookingController.CreateBooking)
					booking.GET("/:id", bookingController.GetBookingByID)
					booking.GET("/:id/receipt", bookingController.DownloadReceipt)
					booking.POST("/:id/payment", bookingController.UploadPaymentProof)
				}
			}

			// Admin routes (require admin role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAdmin(db)) // Apply admin role middleware
			{
				adminRoutes := admin.Group("/routes")
				{
					adminRoutes.GET("", routeController.GetAllRoutes)
				}

				adminRoute := admin.Group("/route")
				{
					adminRoute.POST("", routeController.CreateRoute)
					adminRoute.PUT("/:id", routeController.UpdateRoute)
					adminRoute.DELETE("/:id", routeController.DeleteRoute)
				}

				adminSchedules := admin.Group("/schedules")
				{
					adminSchedules.GET("", scheduleController.GetAllSchedules)
				}

				adminSchedule := admin.Group("/schedule")
				{
					adminSchedule.POST("", scheduleController.CreateSchedule)
					adminSchedule.GET("/:id", scheduleController.GetScheduleByID)
					adminSchedule.PUT("/:id", scheduleController.UpdateSchedule)
					adminSchedule.DELETE("/:id", scheduleController.DeleteSchedule)
				}

				// Admin booking routes
				adminBookings := admin.Group("/bookings")
				{
					adminBookings.GET("", bookingController.GetAllBookings)
				}

				adminBooking := admin.Group("/booking")
				{
					adminBooking.GET("/:id", bookingController.GetBookingByID)
					adminBooking.GET("/:id/payment/download", bookingController.DownloadPaymentProof)
					adminBooking.PUT("/:id/status", bookingController.UpdateBookingStatus)
				}
			}
		}

	}
}
