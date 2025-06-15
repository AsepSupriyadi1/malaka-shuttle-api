package routes

import (
	"malakashuttle/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h *controllers.AuthController) {
	// Auth routes
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}
