package routes

import (
	"malakashuttle/controllers"
	"malakashuttle/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, userController *controllers.UserController) {
	// Admin-only routes for user management
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		// User management routes
		adminRoutes.GET("/users", userController.GetAllUsers)
		adminRoutes.GET("/users/:id", userController.GetUserByID)
		adminRoutes.POST("/users", userController.CreateUser)
		adminRoutes.PUT("/users/:id", userController.UpdateUser)
		adminRoutes.DELETE("/users/:id", userController.DeleteUser)
	}
}
