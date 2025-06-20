package router

import (
	"malakashuttle/controllers"
	"malakashuttle/middleware"
	"malakashuttle/repositories"
	"malakashuttle/routes"
	"malakashuttle/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
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

	// feedback: Ini ntr ganti jadi pake cron job
	bookingScheduler := services.NewBookingScheduler(bookingService)
	bookingScheduler.Start()

	// Apply logging middleware to all API routes
	r.Use(middleware.LoggerMiddleware(), middleware.RequestIDMiddleware())

	// Menginisialisasi grup router untuk API
	router := r.Group("/api")

	// Inisialisasi routes
	routes.TestRoutes(router, testController)
	routes.AuthRoutes(router, authController)
	routes.BookingRoutes(router, bookingController)
	routes.RouteRoutes(router, routeController)
	routes.ScheduleRoutes(router, scheduleController)
}
