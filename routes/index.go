package routes

import (
	controllers "floriangoussin.com/weather-backend/controllers"
	middlewares "floriangoussin.com/weather-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.RouterGroup) {
	initGroup := router.Group("/") 
	initGroup.Use(middlewares.CheckApiKey()) 
	{
		initGroup.POST("/register", controllers.Register)
		initGroup.POST("/login", controllers.Login)
	}
}

func Protected(router *gin.RouterGroup) {
	protectedGroup := router.Group("/")
	protectedGroup.Use(middlewares.CheckApiKey())
	protectedGroup.Use(middlewares.Authenticate())
	{
		// Return city suggestions based on the search terms
		protectedGroup.GET("/cities", controllers.Autocomplete)
	}
}