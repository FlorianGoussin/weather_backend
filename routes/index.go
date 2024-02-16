package routes

import (
	controllers "floriangoussin.com/weather-backend/controllers"
	middlewares "floriangoussin.com/weather-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup) {
	auth := router.Group("/auth") 
	auth.Use(middlewares.CheckApiKey()) 
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}

func Protected(router *gin.RouterGroup) {
	api := router.Group("/")
	api.Use(middlewares.Authenticate())
	{
		// Return city suggestions based on the search terms
		router.GET("/cities", controllers.Autocomplete)
	}
}