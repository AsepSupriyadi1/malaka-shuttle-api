package main

import (
	"log"
	"malakashuttle/config"
	"malakashuttle/entities"
	"malakashuttle/middleware"
	"malakashuttle/repositories"
	"malakashuttle/routes"
	"malakashuttle/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	config.InitLogger()
	logger := config.GetLogger()

	db := config.ConnectDatabase()

	err = entities.AutoMigrate(db)
	if err != nil {
		log.Fatal("Error running auto migration: ", err)
	}
	// Create additional indexes for better performance
	err = entities.CreateIndexes(db)
	if err != nil {
		log.Println("Warning: Error creating indexes: ", err)
	}

	// Initialize booking scheduler
	bookingRepo := repositories.NewBookingRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	bookingService := services.NewBookingService(bookingRepo, scheduleRepo, userRepo)
	bookingScheduler := services.NewBookingScheduler(bookingService)
	bookingScheduler.Start()

	router := gin.New()

	// Add global middleware
	router.Use(middleware.RequestIDMiddleware()) // Add request ID to all requests
	router.Use(gin.Recovery())                   // Add panic recovery
	routes.SetupRoutes(router, db)

	logger.Info("Server starting on port 8080")
	router.Run(":8080")
}
