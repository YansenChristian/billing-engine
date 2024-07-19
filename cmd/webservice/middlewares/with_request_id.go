package middlewares

import (
	"github.com/gin-gonic/gin"
)

func WithRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_id", 1234567)
		c.Next()
	}
}
