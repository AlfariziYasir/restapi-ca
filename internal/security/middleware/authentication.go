package middleware

import (
	"net/http"
	"restapi/internal/security/token"
	"restapi/internal/web"

	"github.com/gin-gonic/gin"
)

func SetupAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := token.TokenValid(c.Request)
		if err != nil {
			web.MarshalError(c, http.StatusUnauthorized, err.Error(), nil)
			c.Abort()
			return
		}

		c.Set("user_id", data.UserId)
		c.Set("username", data.Username)
		c.Set("user_role", data.Role)

		c.Next()
	}
}
