package middleware

import (
	"malakashuttle/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user_role from context (set by AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			err := utils.NewUnauthorizedError("User not authenticated", nil)
			utils.Response.BuildErrorResponse(c, err)
			c.Abort()
			return
		}

		// Check if user has the required role
		if userRole.(string) != requiredRole {
			userEmail, _ := c.Get("user_email")
			forbiddenErr := utils.NewForbiddenErrorWithDetails(
				"Insufficient permissions",
				nil,
				map[string]interface{}{
					"required_role": requiredRole,
					"user_role":     userRole,
					"user_email":    userEmail,
				},
			)
			utils.Response.BuildErrorResponse(c, forbiddenErr)
			c.Abort()
			return
		}

		c.Next()
	}
}
