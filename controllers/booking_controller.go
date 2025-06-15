package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"malakashuttle/dto"
	"malakashuttle/entities"
	"malakashuttle/services"
	"malakashuttle/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookingController struct {
	bookingService *services.BookingService
	validator      *validator.Validate
}

func NewBookingController(bookingService *services.BookingService) *BookingController {
	return &BookingController{
		bookingService: bookingService,
		validator:      validator.New(),
	}
}

// CreateBooking creates a new booking
func (c *BookingController) CreateBooking(ctx *gin.Context) {
	// Get user email from JWT token
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.CreateBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := c.validator.Struct(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}
	// Create booking
	booking, err := c.bookingService.CreateBooking(userEmail.(string), req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "already booked") {
			utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "cannot book past") {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create booking", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Booking created successfully", booking)
}

// GetBookingByID gets booking by ID
func (c *BookingController) GetBookingByID(ctx *gin.Context) {
	// Get booking ID from URL
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID", nil)
		return
	}
	// Get user email from JWT token
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Check if user is staff (can view all bookings) or regular user (can only view own bookings)
	userRole, _ := ctx.Get("user_role")
	var userEmailPtr *string
	if userRole != "staff" && userRole != "admin" {
		email := userEmail.(string)
		userEmailPtr = &email
	}

	booking, err := c.bookingService.GetBookingDetailByID(uint(bookingID), userEmailPtr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get booking", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Booking retrieved successfully", booking)
}

// GetUserBookings gets all bookings for the authenticated user
func (c *BookingController) GetUserBookings(ctx *gin.Context) {
	// Get user email from JWT token
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get pagination parameters
	params := utils.GetPaginationParams(ctx)
	statusStr := ctx.Query("status")

	// Parse status filter
	var statusFilter []entities.BookingStatus
	if statusStr != "" {
		statuses := strings.Split(statusStr, ",")
		for _, s := range statuses {
			statusFilter = append(statusFilter, entities.BookingStatus(strings.TrimSpace(s)))
		}
	}

	bookings, err := c.bookingService.GetUserBookingsList(userEmail.(string), params, statusFilter)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get bookings", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Bookings retrieved successfully", bookings)
}

// GetAllBookings gets all bookings (for staff only)
func (c *BookingController) GetAllBookings(ctx *gin.Context) {
	// Get pagination parameters
	params := utils.GetPaginationParams(ctx)
	statusStr := ctx.Query("status")

	// Parse status filter
	var statusFilter []entities.BookingStatus
	if statusStr != "" {
		statuses := strings.Split(statusStr, ",")
		for _, s := range statuses {
			statusFilter = append(statusFilter, entities.BookingStatus(strings.TrimSpace(s)))
		}
	}

	bookings, err := c.bookingService.GetAllBookingsList(params, statusFilter)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get bookings", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Bookings retrieved successfully", bookings)
}

// UploadPaymentProof uploads payment proof for a booking
func (c *BookingController) UploadPaymentProof(ctx *gin.Context) {
	// Get booking ID from URL
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID", nil)
		return
	}

	// Get user email from JWT token
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get payment method from form
	paymentMethod := ctx.PostForm("payment_method")
	if paymentMethod == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Payment method is required", nil)
		return
	}

	// Get uploaded file
	file, err := ctx.FormFile("proof_image")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Payment proof image is required", err.Error())
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "File size too large (max 5MB)", nil)
		return
	}

	// Validate file type
	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png"}
	fileHeader := file.Header.Get("Content-Type")
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if fileHeader == allowedType {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid file type. Only JPEG, JPG, and PNG are allowed", nil)
		return
	}

	// Upload payment proof
	err = c.bookingService.UploadPaymentProof(uint(bookingID), userEmail.(string), file, paymentMethod)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not eligible") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "expired") {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "already uploaded") {
			utils.ErrorResponse(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to upload payment proof", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Payment proof uploaded successfully", nil)
}

// UpdateBookingStatus updates booking status (for staff only)
func (c *BookingController) UpdateBookingStatus(ctx *gin.Context) {
	// Get booking ID from URL
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID", nil)
		return
	}

	var req dto.UpdateBookingStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validate request
	if err := c.validator.Struct(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Update booking status
	err = c.bookingService.UpdateBookingStatus(uint(bookingID), req.Status)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "invalid status") || strings.Contains(err.Error(), "not in waiting") {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update booking status", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Booking status updated successfully", nil)
}

// GetAvailableSeats gets available seats for a schedule
func (c *BookingController) GetAvailableSeats(ctx *gin.Context) {
	// Get schedule ID from URL
	scheduleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid schedule ID", nil)
		return
	}

	seats, err := c.bookingService.GetAvailableSeats(uint(scheduleID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get available seats", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Available seats retrieved successfully", seats)
}

// DownloadPaymentProof downloads payment proof file (for staff only)
func (c *BookingController) DownloadPaymentProof(ctx *gin.Context) {
	// Get booking ID from URL
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID", nil)
		return
	}
	// Get booking to check if payment proof exists
	_, err = c.bookingService.GetBookingByID(uint(bookingID), nil) // nil = staff can access any booking
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get booking", err.Error())
		return
	}

	// Check if payment proof exists
	payment, err := c.bookingService.GetPaymentByBookingID(uint(bookingID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Payment proof not found", nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get payment", err.Error())
		return
	}

	if payment.ProofImageURL == "" {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Payment proof file not found", nil)
		return
	}

	// Construct file path
	filePath := "." + payment.ProofImageURL // Remove leading slash and add current directory

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Payment proof file not found on disk", nil)
		return
	}

	// Set headers for file download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=payment_proof_%d%s", bookingID, filepath.Ext(filePath)))
	ctx.Header("Content-Type", "application/octet-stream")

	// Serve the file
	ctx.File(filePath)
}

// DownloadReceipt generates and downloads booking receipt (for users)
func (c *BookingController) DownloadReceipt(ctx *gin.Context) {
	// Get booking ID from URL
	bookingID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid booking ID", nil)
		return
	}
	// Get user email from JWT token
	userEmail, exists := ctx.Get("user_email")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Check if user is staff (can access any receipt) or regular user (can only access own receipt)
	userRole, _ := ctx.Get("user_role")
	var userEmailPtr *string
	if userRole != "staff" && userRole != "admin" {
		email := userEmail.(string)
		userEmailPtr = &email
	}

	// Generate receipt
	receiptPath, err := c.bookingService.GenerateBookingReceipt(uint(bookingID), userEmailPtr)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.ErrorResponse(ctx, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "only be generated for successful") {
			utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate receipt", err.Error())
		return
	}

	// Check if file exists
	if _, err := os.Stat(receiptPath); os.IsNotExist(err) {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Receipt file not found", nil)
		return
	}

	// Set headers for PDF download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=malaka_shuttle_receipt_%d.pdf", bookingID))
	ctx.Header("Content-Type", "application/pdf")

	// Serve the file
	ctx.File(receiptPath)
}
