package routes

import "github.com/gin-gonic/gin"
controllers "floriangoussin.com/weather-backend/controllers"

func InitializeRoutes(router *gin.Engine) {
	router.POST("api/register", controllers.Register)
  router.POST("api/login", controllers.Login)
	
	// Return city suggestions based on the search terms
  router.GET("api/cities", handleAutocomplete)
}