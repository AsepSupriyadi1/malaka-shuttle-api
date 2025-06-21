package routes

import (
	"malakashuttle/constants"
	"malakashuttle/controllers"
	"malakashuttle/middleware"

	"github.com/gin-gonic/gin"
)

func RouteRoutes(r *gin.RouterGroup, h *controllers.RouteController) {

	adminRoutes := r.Group("admin/routes")
	adminRoutes.Use(middleware.AuthMiddleware(), middleware.RequireRole(constants.ROLE_ADMIN))
	adminRoutes.GET("", h.GetAllRoutes)
	adminRoutes.GET("/:id", h.GetRouteByID)
	adminRoutes.POST("", h.CreateRoute)
	adminRoutes.PUT("/:id", h.UpdateRoute)
	adminRoutes.DELETE("/:id", h.DeleteRoute)

}
