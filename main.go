package main

import (
	"fmt"
	"log"
	"malakashuttle/config"
	"malakashuttle/middleware"
	myRouter "malakashuttle/router"
	"os"

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

	err = config.ResetDatabase(db)
	if err != nil {
		log.Fatal("Error resetting database: ", err)
	}

	router := gin.New()

	// Add global middleware
	router.Use(middleware.RequestIDMiddleware()) // Add request ID to all requests
	myRouter.InitRoutes(router, db)

	logger.Info("Server starting on port 8080")
	router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
