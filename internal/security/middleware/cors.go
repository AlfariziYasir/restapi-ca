package middleware

import (
	"reflect"
	"restapi/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	var accessControlAllowHeaders = strings.Join(
		[]string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"X-XSRF-TOKEN",
		},
		", ",
	)

	return func(c *gin.Context) {
		http_origin := c.GetHeader("origin")
		whitelistHost := strings.Split(config.Cfg().WhitelistHost, ",")
		exist, _ := indexOf(http_origin, whitelistHost)
		if exist {
			c.Writer.Header().Set("Access-Control-Allow-Origin", http_origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, UPDATE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", accessControlAllowHeaders)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Referrer-Policy", "same-origin")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=15724800; includeSubDomains")

		c.Writer.Header().Set("Access-Control-Max-Age", "0")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func indexOf(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
