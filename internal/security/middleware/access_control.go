package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		role := c.MustGet("user_role").(string)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
				"error":   http.StatusInternalServerError,
				"message": "user unauthorized",
				"data":    nil,
			})
		}

		c.Next()
	}
}
