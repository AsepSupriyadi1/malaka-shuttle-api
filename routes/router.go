package routes

import (
	"malakashuttle/controllers"
	"malakashuttle/middleware"
	"malakashuttle/repositories"
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

	bookingScheduler := services.NewBookingScheduler(bookingService)
	bookingScheduler.Start()

	// Apply logging middleware to all API routes
	router := r.Group("/api")
	router.GET("/ping", testController.Ping)
	router.Use(middleware.LoggerMiddleware())
	AuthRoutes(router, authController)
	BookingRoutes(router, bookingController)
	RouteRoutes(router, routeController)
	ScheduleRoutes(router, scheduleController)
}
