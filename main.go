package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"os/signal"
	"syscall"

	database "floriangoussin.com/weather-backend/database"
	routes "floriangoussin.com/weather-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	docs "floriangoussin.com/weather-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
  if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
  // Get engine
  router := gin.New()
  router.Use(gin.Logger())
  router.ForwardedByClientIP = true
  router.SetTrustedProxies([]string{"127.0.0.1"})

  // programmatically set swagger info
	docs.SwaggerInfo.Title = "Weather mobile app API"
	docs.SwaggerInfo.Description = "Interface between the app and the Weather API. More info at: https://www.weatherapi.com/."
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = os.Getenv("APP_URI")
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

  // log.Println("GOROOT:", runtime.GOROOT())
  // use ginSwagger middleware to serve the API docs
  router.GET("/", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
  })

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

  api := router.Group("/api/v1")
  routes.Init(api)
  routes.Protected(api)

  appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}
  // router.Run(":" + appPort)
  // Create a server instance with a timeout
	server := &http.Server{
		Addr: ":" + appPort,
		Handler: router,
	}

  // InitializeDatabase will create Weather database and
	// the Cities collection with preloaded data
  database.Connect()
  database.Initialize(database.Client)

  go func() {
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server failed to start: %v", err)
    }
  }()

  // Handle shutdown signals
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
  <-sigs

  // Close db connection
  database.Disconnect()

  // Shutdown the server gracefully
  if err := server.Shutdown(context.Background()); err != nil {
    log.Printf("server shutdown error: %v", err)
  }

}