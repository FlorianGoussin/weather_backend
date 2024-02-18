package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

  api := router.Group("/api/v1")
  api.GET("/", func (c *gin.Context)  {
    responseData := gin.H{
      "message": "Success",
    }
    // Return the JSON response with status code 200
    c.JSON(http.StatusOK, responseData)
  })
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

  // Handle graceful shutdown
  go func() {
    quit := make(chan os.Signal, 1) // Create channel
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down gracefully")

    // Create a context with a timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Close db connection
    database.Disconnect()

    // Shutdown the server
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown failed: %v", err)
    }

    log.Println("Server gracefully stopped")
  }()
}