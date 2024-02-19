package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"os/signal"
	"syscall"

	// "time"

	database "floriangoussin.com/weather-backend/database"
	routes "floriangoussin.com/weather-backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
  router.GET("/", func (c *gin.Context)  {
    responseData := gin.H{
      "message": "Root route test Heroku",
    }
    // Return the JSON response with status code 200
    c.JSON(http.StatusOK, responseData)
  })

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