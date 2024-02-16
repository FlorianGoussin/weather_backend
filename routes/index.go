package routes

import (
	controllers "floriangoussin.com/weather-backend/controllers"
	middlewares "floriangoussin.com/weather-backend/middleware"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.POST("api/register", controllers.Register)
  router.POST("api/login", controllers.Login)
}

func Protected(router *gin.Engine) {
	router.Use(middlewares.Authenticate)

	// Return city suggestions based on the search terms
  router.GET("api/cities", controllers.Autocomplete)
}