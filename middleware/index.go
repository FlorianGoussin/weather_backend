package middlewares

import (
	"net/http"

	helper "floriangoussin.com/weather-backend/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
				c.Abort()
				return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				c.Abort()
				return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		// c.Set("user_type", claims.User_type)
		c.Next()
	}
}

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