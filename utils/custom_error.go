package utils

import (
	"fmt"
	"log"
)

type CustomError struct {
	StatusCode int
	Message    string
	Err        error
	Details    interface{} // Untuk menyimpan detail data tambahan
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		log.Printf("StatusCode: %d, Message: %s, Detail: %v", e.StatusCode, e.Message, e.Err)
		return fmt.Sprintf("%v", e.Err)
	}
	return ""
}

func (e *CustomError) Unwrap() error {
	return e.Err
}

func NewCustomError(statusCode int, message string, err error) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
		Details:    nil,
	}
}

// NewCustomErrorWithDetails membuat custom error dengan detail data tambahan
func NewCustomErrorWithDetails(statusCode int, message string, err error, details interface{}) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
		Details:    details,
	}
}

// Helper functions untuk membuat custom error dengan mudah
func NewBadRequestError(message string, err error) *CustomError {
	return NewCustomError(400, message, err)
}

func NewUnauthorizedError(message string, err error) *CustomError {
	return NewCustomError(401, message, err)
}

func NewForbiddenError(message string, err error) *CustomError {
	return NewCustomError(403, message, err)
}

func NewNotFoundError(message string, err error) *CustomError {
	return NewCustomError(404, message, err)
}

func NewConflictError(message string, err error) *CustomError {
	return NewCustomError(409, message, err)
}

func NewValidationError(message string, err error) *CustomError {
	return NewCustomError(422, message, err)
}

func NewInternalServerError(message string, err error) *CustomError {
	return NewCustomError(500, message, err)
}

// Helper functions dengan detail data tambahan
func NewBadRequestErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(400, message, err, details)
}

func NewUnauthorizedErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(401, message, err, details)
}

func NewForbiddenErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(403, message, err, details)
}

func NewNotFoundErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(404, message, err, details)
}

func NewConflictErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(409, message, err, details)
}

func NewValidationErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(422, message, err, details)
}

func NewInternalServerErrorWithDetails(message string, err error, details interface{}) *CustomError {
	return NewCustomErrorWithDetails(500, message, err, details)
}
