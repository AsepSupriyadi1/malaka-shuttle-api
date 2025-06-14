package middleware

import (
	"malakashuttle/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Membaca header Authorization dari request
		authHeader := c.GetHeader("Authorization")

		// Memeriksa apakah header Authorization ada dan tidak kosong
		if authHeader == "" {
			err := utils.NewUnauthorizedError("Authorization header is required", nil)
			utils.Response.BuildErrorResponse(c, err)
			c.Abort()
			return
		}

		// Memeriksa apakah token diawali dengan "Bearer "
		parts := strings.Split(authHeader, " ")

		// Memeriksa panjang parts dan format Bearer token
		if len(parts) != 2 || parts[0] != "Bearer" {
			err := utils.NewUnauthorizedErrorWithDetails(
				"Invalid authorization header format",
				nil,
				map[string]string{
					"expected_format": "Bearer <token>",
					"received_format": authHeader,
				},
			)
			utils.Response.BuildErrorResponse(c, err)
			c.Abort()
			return
		}

		// Validasi token
		userId, err := utils.ValidateToken(parts[1])

		if err != nil {
			tokenErr := utils.NewUnauthorizedErrorWithDetails(
				"Invalid or expired token",
				err,
				map[string]string{
					"token_error": err.Error(),
				},
			)
			utils.Response.BuildErrorResponse(c, tokenErr)
			c.Abort()
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}
