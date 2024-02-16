package middlewares

import (
	"net/http"
	"os"

	helper "floriangoussin.com/weather-backend/helpers"
	"github.com/gin-gonic/gin"
)

var clientApiKey = os.Getenv("CLIENT_API_KEY")
var validAPIKeys = map[string]bool{
	clientApiKey: true,
}

func CheckApiKey() gin.HandlerFunc { return _checkApiKey }
func _checkApiKey(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	// Check if API key is provided
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is missing"})
		c.Abort()
		return
	}

	// Check if the API key is valid
	if !validAPIKeys[apiKey] {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		c.Abort()
		return
	}
}

func Authenticate() gin.HandlerFunc { return _authenticate }
func _authenticate(c *gin.Context) {
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
	c.Next()
}

// TODO: See if of any use
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