package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Log the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current working directory: %s", dir)

	// Check if the binary executable exists
	if _, err := os.Stat("./weather-backend"); err == nil {
		log.Println("Executable binary found.")
	} else if os.IsNotExist(err) {
		log.Println("Executable binary not found.")
	} else {
		log.Fatalf("Error checking for executable binary: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Define a route handler for the root path
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Run the server on port 8080
	router.Run(":8080")
}