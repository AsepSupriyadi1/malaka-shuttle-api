package middleware

import (
	"malakashuttle/repositories"
	"malakashuttle/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(db *gorm.DB, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_id from context (set by AuthMiddleware)
		userID, exists := c.Get("user_id")
		if !exists {
			err := utils.NewUnauthorizedError("User not authenticated", nil)
			utils.Response.BuildErrorResponse(c, err)
			c.Abort()
			return
		}

		// Get user from database to check role
		userRepo := repositories.NewUserRepository(db)
		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			notFoundErr := utils.NewNotFoundErrorWithDetails(
				"User not found",
				err,
				map[string]interface{}{
					"user_id": userID,
				},
			)
			utils.Response.BuildErrorResponse(c, notFoundErr)
			c.Abort()
			return
		}

		// Check if user has the required role
		if user.Role != requiredRole {
			forbiddenErr := utils.NewForbiddenErrorWithDetails(
				"Insufficient permissions",
				nil,
				map[string]interface{}{
					"required_role": requiredRole,
					"user_role":     user.Role,
					"user_id":       userID,
				},
			)
			utils.Response.BuildErrorResponse(c, forbiddenErr)
			c.Abort()
			return
		}

		// Set user role in context for later use
		c.Set("user_role", user.Role)
		c.Next()
	}
}

// RequireAdmin is a convenience function for admin role requirement
func RequireAdmin(db *gorm.DB) gin.HandlerFunc {
	return RequireRole(db, "admin")
}

func RequireUser(db *gorm.DB) gin.HandlerFunc {
	return RequireRole(db, "user")
}

func RequireStaff(db *gorm.DB) gin.HandlerFunc {
	return RequireRole(db, "staff")
}
