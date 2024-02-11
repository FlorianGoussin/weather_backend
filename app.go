package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "floriangoussin.com/weather-backend/database"
  "github.com/joho/godotenv"
  "log"
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "strings"
  // "fmt"
)

func main() {
  if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

  r := gin.Default()

  // r.GET("/", func(c *gin.Context) {
  //   c.JSON(http.StatusOK, gin.H{"data": "Root route"})    
  // })

  // Return city suggestions based on the search terms
  r.GET("/cities", handleAutocomplete)

  r.Run()
}

func handleAutocomplete(c *gin.Context) {
  weatherDatabase := database.Connect()

  searchTerm := c.Query("searchTerm")

  // c.JSON(http.StatusOK, gin.H{"data": "Cities autocomplete"})    
  collection := weatherDatabase.Collection("Cities")

  // Fetch all documents from the collection
  cursor, err := collection.Find(context.Background(), bson.M{})
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  defer cursor.Close(context.Background())
  
  // Check if any documents were found
  var suggestions []bson.M
  for cursor.Next(context.Background()) {
    var result bson.M
    err := cursor.Decode(&result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

    // Check if the value stored in the "name" field is actually a string
    name, ok := result["name"].(string)
    if !ok {
      // Handle the case where the "name" field is not a string
      // var typeName string
      // if result["name"] != nil {
      //     typeName = fmt.Sprintf("%T", result["name"])
      // } else {
      //     typeName = "nil"
      // }
      // log.Printf("Name field (%s) is not a string\n", typeName)
      // log.Printf("Result record: %+v\n", result)
      continue // Skip this iteration and move to the next document
    }
    if strings.HasPrefix(strings.ToLower(name), strings.ToLower(searchTerm)) {
			suggestions = append(suggestions, result)
		}
  }
  defer database.Disconnect()

  // Return the suggestions
  c.JSON(http.StatusOK, gin.H{
    "message": "Cities suggestions",
    "searchTerm": searchTerm,
    "suggestions": suggestions,
  })
}