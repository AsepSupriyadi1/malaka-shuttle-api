package controllers

import (
	"malakashuttle/dto"
	"malakashuttle/services"
	"malakashuttle/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RouteController struct {
	routeService services.RouteService
}

func NewRouteController(routeService services.RouteService) *RouteController {
	return &RouteController{
		routeService: routeService,
	}
}

func (rc *RouteController) CreateRoute(c *gin.Context) {
	var req dto.RouteRequest

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response.HandleValidationError(c, err, req)
		return
	}

	// Call service
	response, err := rc.routeService.CreateRoute(req)
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.Created(c, "Route created successfully", response)
}

func (rc *RouteController) GetRouteByID(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Response.BadRequest(c, "Invalid route ID", nil)
		return
	}

	// Call service
	response, err := rc.routeService.GetRouteByID(uint(id))
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.OK(c, "Route retrieved successfully", response)
}

func (rc *RouteController) UpdateRoute(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Response.BadRequest(c, "Invalid route ID", nil)
		return
	}

	var req dto.RouteUpdateRequest

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Response.HandleValidationError(c, err, req)
		return
	}

	// Call service
	response, err := rc.routeService.UpdateRoute(uint(id), req)
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.OK(c, "Route updated successfully", response)
}

func (rc *RouteController) DeleteRoute(c *gin.Context) {
	// Get ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.Response.BadRequest(c, "Invalid route ID", nil)
		return
	}

	// Call service
	err = rc.routeService.DeleteRoute(uint(id))
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.OK(c, "Route deleted successfully", nil)
}

func (rc *RouteController) GetAllRoutes(c *gin.Context) {
	// Get pagination parameters
	params := utils.GetPaginationParams(c)

	// Call service
	response, err := rc.routeService.GetAllRoutes(params)
	if err != nil {
		utils.Response.BuildErrorResponse(c, err)
		return
	}

	utils.Response.OK(c, "Routes retrieved successfully", response)
}
