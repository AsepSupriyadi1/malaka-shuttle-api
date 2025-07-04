package routes

import (
	"malakashuttle/constants"
	"malakashuttle/controllers"
	"malakashuttle/middleware"

	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(r *gin.RouterGroup, h *controllers.ScheduleController) {
	userRoutes := r.Group("/schedules")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_USER))
	userRoutes.GET("/search", h.SearchSchedules)
	userRoutes.GET("/:id", h.GetScheduleByID)

	adminRoutes := r.Group("admin/schedules")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_ADMIN))
	adminRoutes.GET("", h.GetAllSchedules)
	adminRoutes.POST("", h.CreateSchedule)
	adminRoutes.GET("/:id", h.GetScheduleByID)
	adminRoutes.PUT("/:id", h.UpdateSchedule)
	adminRoutes.DELETE("/:id", h.DeleteSchedule)
}
