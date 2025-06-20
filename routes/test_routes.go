package routes

import (
	"malakashuttle/controllers"

	"github.com/gin-gonic/gin"
)

func TestRoutes(r *gin.RouterGroup, h *controllers.TestController) {
	r.GET("/ping", h.Ping)

}
