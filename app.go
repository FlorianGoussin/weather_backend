package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"floriangoussin.com/weather-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
  if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
  // InitializeDatabase will create Weather database and
	// the Cities collection with preloaded data
  client := database.Connect()
  database.Initialize(client)
  database.Disconnect(client)

  // Get engine
  engine := gin.Default()

  // Return city suggestions based on the search terms
  engine.GET("/cities", handleAutocomplete)

  engine.Run()
}

func handleAutocomplete(c *gin.Context) {
  client := database.Connect()

  // Open connection
  collection := client.Database("Weather").Collection("Cities")

  searchTerm := c.Query("searchTerm")
  log.Println("handleAutocomplete searchTerm", searchTerm)

  // Define a filter to search for suggestions based on the searchTerm in the "name" field
  pattern := fmt.Sprintf("^%s", searchTerm)
  filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}}

  // Execute the find operation
  cursor, err := collection.Find(context.Background(), filter)
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  defer cursor.Close(context.Background())

  // Decode results
  log.Println("Number of results:", cursor.RemainingBatchLength())
	// var results []bson.M
	var results []database.City
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
    return
	}

  // Close connection
  defer database.Disconnect(client)

  // Return the suggestions
  c.JSON(http.StatusOK, gin.H{
    "message": "Cities suggestions",
    "searchTerm": searchTerm,
    "suggestions": results,
  })
}