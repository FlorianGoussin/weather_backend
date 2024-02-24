package routes

import (
	controllers "floriangoussin.com/weather-backend/controllers"
	middlewares "floriangoussin.com/weather-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Init(api *gin.RouterGroup) {
	initGroup := api.Group("/") 
	initGroup.Use(middlewares.CheckApiKey()) 
	{
		initGroup.POST("/register", controllers.Register)
		initGroup.POST("/login", controllers.Login)
	}
}

func Protected(api *gin.RouterGroup) {
	protectedGroup := api.Group("/")
	protectedGroup.Use(middlewares.CheckApiKey())
	protectedGroup.Use(middlewares.Authenticate())
	{
		// Return city suggestions based on the search terms
		protectedGroup.GET("/cities", controllers.Autocomplete)

		// Current weather by location
		// protectedGroup.GET("/currentWeather", controllers.GetAllCurrentWeatherByLocation())
		protectedGroup.POST("/currentWeather", controllers.AddCurrentWeatherByCity())
		// protectedGroup.DELETE("/currentWeather", controllers.RemoveCurrentWeatherByLocation())
	}
}