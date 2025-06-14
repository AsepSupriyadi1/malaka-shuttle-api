package utils

import (
	"errors"
	"malakashuttle/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONResponse adalah struktur response yang konsisten
type JSONResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ResponseHelper berisi helper functions untuk response
type ResponseHelper struct{}

// NewResponseHelper membuat instance baru dari ResponseHelper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// BuildSuccessResponse membuat response sukses dengan format yang konsisten
func (r *ResponseHelper) BuildSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := JSONResponse{
		Code:    constants.RESPONSE_SUCCESS,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// BuildErrorResponse membuat response error berdasarkan custom error atau error biasa
func (r *ResponseHelper) BuildErrorResponse(c *gin.Context, err error) {
	var customErr *CustomError
	var statusCode int
	var message string
	var responseData interface{}

	if errors.As(err, &customErr) {
		statusCode = customErr.StatusCode
		message = customErr.Message

		// Prioritas detail data: Details -> Err -> nil
		if customErr.Details != nil {
			responseData = customErr.Details
		} else if customErr.Err != nil {
			responseData = map[string]interface{}{
				"error": customErr.Err.Error(),
			}
		}
	} else {
		statusCode = http.StatusInternalServerError
		message = "An unexpected error occurred"
		responseData = map[string]interface{}{
			"error": err.Error(),
		}
	}

	// Map status code ke response code
	var code string
	switch statusCode {
	case http.StatusBadRequest:
		code = constants.RESPONSE_BAD_REQUEST
	case http.StatusNotFound:
		code = constants.RESPONSE_NOT_FOUND
	case http.StatusUnauthorized:
		code = constants.RESPONSE_UNAUTHENTICATED
	case http.StatusForbidden:
		code = constants.RESPONSE_FORBIDDEN
	case http.StatusTooManyRequests:
		code = constants.RESPONSE_TOO_MANY_REQUESTS
	case http.StatusConflict:
		code = constants.RESPONSE_CONFLICT
	case http.StatusUnprocessableEntity:
		code = constants.RESPONSE_VALIDATION_ERROR
	case http.StatusInternalServerError:
		code = constants.RESPONSE_INTERNAL_SERVER_ERROR
	default:
		code = constants.RESPONSE_INTERNAL_SERVER_ERROR
	}

	response := JSONResponse{
		Code:    code,
		Message: message,
		Data:    responseData,
	}

	c.JSON(statusCode, response)
}

// Helper functions untuk response umum

// OK - 200 Success
func (r *ResponseHelper) OK(c *gin.Context, message string, data interface{}) {
	r.BuildSuccessResponse(c, http.StatusOK, message, data)
}

// Created - 201 Created
func (r *ResponseHelper) Created(c *gin.Context, message string, data interface{}) {
	r.BuildSuccessResponse(c, http.StatusCreated, message, data)
}

// BadRequest - 400 Bad Request
func (r *ResponseHelper) BadRequest(c *gin.Context, message string, data interface{}) {
	err := NewBadRequestErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// Unauthorized - 401 Unauthorized
func (r *ResponseHelper) Unauthorized(c *gin.Context, message string, data interface{}) {
	err := NewUnauthorizedErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// Forbidden - 403 Forbidden
func (r *ResponseHelper) Forbidden(c *gin.Context, message string, data interface{}) {
	err := NewForbiddenErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// NotFound - 404 Not Found
func (r *ResponseHelper) NotFound(c *gin.Context, message string, data interface{}) {
	err := NewNotFoundErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// Conflict - 409 Conflict
func (r *ResponseHelper) Conflict(c *gin.Context, message string, data interface{}) {
	err := NewConflictErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// ValidationError - 422 Unprocessable Entity
func (r *ResponseHelper) ValidationError(c *gin.Context, message string, data interface{}) {
	err := NewValidationErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// InternalServerError - 500 Internal Server Error
func (r *ResponseHelper) InternalServerError(c *gin.Context, message string, data interface{}) {
	err := NewInternalServerErrorWithDetails(message, nil, data)
	r.BuildErrorResponse(c, err)
}

// HandleValidationError adalah helper khusus untuk validation error
// Mengembalikan data yang diinputkan client untuk memudahkan debugging
func (r *ResponseHelper) HandleValidationError(c *gin.Context, err error, inputData interface{}) {
	r.ValidationError(c, "Validation failed", inputData)
}

// Global helper instance
var Response = NewResponseHelper()

// Simplified helper functions for easier usage
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	Response.BuildSuccessResponse(c, statusCode, message, data)
}

func ErrorResponse(c *gin.Context, statusCode int, message string, details interface{}) {
	var err error
	switch statusCode {
	case http.StatusBadRequest:
		err = NewBadRequestErrorWithDetails(message, nil, details)
	case http.StatusUnauthorized:
		err = NewUnauthorizedErrorWithDetails(message, nil, details)
	case http.StatusForbidden:
		err = NewForbiddenErrorWithDetails(message, nil, details)
	case http.StatusNotFound:
		err = NewNotFoundErrorWithDetails(message, nil, details)
	case http.StatusConflict:
		err = NewConflictErrorWithDetails(message, nil, details)
	case http.StatusUnprocessableEntity:
		err = NewValidationErrorWithDetails(message, nil, details)
	default:
		err = NewInternalServerErrorWithDetails(message, nil, details)
	}
	Response.BuildErrorResponse(c, err)
}
