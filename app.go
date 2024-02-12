package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"floriangoussin.com/weather-backend/database"
	// "floriangoussin.com/weather-backend/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// var IPWhitelist = map[string]bool{
//   "127.0.0.1": true,
//   "111.2.3.4": true,
//   "::1": true,
// }

func main() {
  if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
  r := gin.Default()

  // r.ForwardedByClientIP = true
  // r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.2", "10.0.0.0/8"})

  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Root endpoint works in Docker",
    })
  })

  // Return city suggestions based on the search terms
  r.GET("/cities", handleAutocomplete)

  // restrictedPage := r.Group("/")
  // restrictedPage.Use(middlewares.IPWhiteListMiddleware(IPWhitelist))
  // restrictedPage.GET("/adminZone", func(c *gin.Context) {
  //   c.JSON(http.StatusOK, gin.H{
  //       "message": "This endpoint is secured with IP whitelisting!",
  //   })
  // })

  r.Run()
}

func handleAutocomplete(c *gin.Context) {
  weatherDatabase := database.Connect()

  searchTerm := c.Query("searchTerm")
  log.Println("handleAutocomplete searchTerm", searchTerm)

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
    name, _ := result["name"].(string)
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