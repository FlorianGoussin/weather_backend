package routes

import (
	controllers "floriangoussin.com/weather-backend/controllers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.POST("api/register", controllers.Register)
  router.POST("api/login", controllers.Login)
	
	// Return city suggestions based on the search terms
  router.GET("api/cities", controllers.Autocomplete)
}