package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhiteListMiddleware(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIP := c.ClientIP()

		if !whitelist[userIP] {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"error": "You are not authorized to access this resource!",
				})
		} else {
				c.Next()
		}
	}
}