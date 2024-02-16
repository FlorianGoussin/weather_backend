package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	database "floriangoussin.com/weather-backend/database"
	models "floriangoussin.com/weather-backend/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Autocomplete(c *gin.Context) {
	client := database.Client
  citiesCollection := client.Database("Weather").Collection("Cities")

  searchTerm := c.Query("searchTerm")
  log.Println("handleAutocomplete searchTerm", searchTerm)

  // TODO: Sanitize searchTerm

  // Define a filter to search for suggestions based on the searchTerm in the "name" field
  pattern := fmt.Sprintf("^%s", searchTerm)
  filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}}

  // Execute the find operation
  cursor, err := citiesCollection.Find(context.Background(), filter)
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  defer cursor.Close(context.Background())

  // Decode results
  log.Println("Number of results:", cursor.RemainingBatchLength())
	// var results []bson.M
	var results []models.City
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
    return
	}

  // Return the suggestions
  c.JSON(http.StatusOK, gin.H{
    "message": "Cities suggestions",
    "searchTerm": searchTerm,
    "suggestions": results,
  })
}