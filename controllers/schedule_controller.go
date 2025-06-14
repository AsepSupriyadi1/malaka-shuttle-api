package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"malakashuttle/dto"
	"malakashuttle/services"
	"malakashuttle/utils"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	scheduleService *services.ScheduleService
}

func NewScheduleController(scheduleService *services.ScheduleService) *ScheduleController {
	return &ScheduleController{
		scheduleService: scheduleService,
	}
}

// CreateSchedule - Create new schedule (Admin only)
func (c *ScheduleController) CreateSchedule(ctx *gin.Context) {
	var req dto.CreateScheduleRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request data", req)
		return
	}

	schedule, err := c.scheduleService.CreateSchedule(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create schedule", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Schedule created successfully", schedule)
}

// UpdateSchedule - Update schedule (Admin only)
func (c *ScheduleController) UpdateSchedule(ctx *gin.Context) {
	// Get schedule ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid schedule ID", nil)
		return
	}

	var req dto.UpdateScheduleRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Update schedule
	schedule, err := c.scheduleService.UpdateSchedule(uint(id), req)
	if err != nil {
		if err.Error() == "schedule not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Schedule not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Schedule updated successfully", schedule)
}

// DeleteSchedule - Delete schedule (Admin only)
func (c *ScheduleController) DeleteSchedule(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid schedule ID", nil)
		return
	}

	if err := c.scheduleService.DeleteSchedule(uint(id)); err != nil {
		if err.Error() == "schedule not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Schedule not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Schedule deleted successfully", nil)
}

// GetAllSchedules - Get all schedules with pagination (Admin only)
func (c *ScheduleController) GetAllSchedules(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	schedules, err := c.scheduleService.GetAllSchedules(page, limit)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get schedules", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Schedules retrieved successfully", schedules)
}

// SearchSchedules - Search schedules by origin, destination, and departure date
func (c *ScheduleController) SearchSchedules(ctx *gin.Context) {
	var req dto.ScheduleSearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	schedules, err := c.scheduleService.SearchSchedules(req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to search schedules", err.Error())
		return
	}

	if len(schedules.Results) == 0 {
		utils.SuccessResponse(ctx, http.StatusOK, "No schedules found for the given criteria", schedules)
		return
	}

	message := fmt.Sprintf("Found %d schedule(s)", len(schedules.Results))
	utils.SuccessResponse(ctx, http.StatusOK, message, schedules)
}

// GetScheduleByID - Get schedule by ID (Admin)
func (c *ScheduleController) GetScheduleByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid schedule ID", nil)
		return
	}

	schedule, err := c.scheduleService.GetScheduleByID(uint(id), true) // true untuk admin
	if err != nil {
		if err.Error() == "schedule not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Schedule not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Schedule retrieved successfully", schedule)
}

// GetScheduleByIDForUser - Get schedule by ID (User)
func (c *ScheduleController) GetScheduleByIDForUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid schedule ID", nil)
		return
	}

	schedule, err := c.scheduleService.GetScheduleByID(uint(id), false) // false untuk user
	if err != nil {
		if err.Error() == "schedule not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Schedule not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Schedule retrieved successfully", schedule)
}
