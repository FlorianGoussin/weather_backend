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

// var IPWhitelist = map[string]bool{
//   "127.0.0.1": true,
//   "111.2.3.4": true,
//   "::1": true,
// }

type City struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

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

  // Open connection
  collection := weatherDatabase.Collection("Cities")

  searchTerm := c.Query("searchTerm")
  log.Println("handleAutocomplete searchTerm", searchTerm)

  // Define a filter to search for suggestions based on the searchTerm in the "name" field
  pattern := fmt.Sprintf("^%s", searchTerm)
  filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}}
  // filter := bson.D{{Key: "name", Value: primitive.Regex{Pattern: pattern, Options: ""}}}

  // Execute the find operation
  cursor, err := collection.Find(context.Background(), filter)
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  defer cursor.Close(context.Background())

  // Construct a regular expression pattern for the string variable
	// pattern := fmt.Sprintf("^%s", searchTerm)
	// regex := primitive.Regex{Pattern: pattern, Options: ""}
  // log.Println("regex", regex)

  // // Stages and pipeline:
  // matchStage := bson.D{
  //   {Key: "$match", Value: bson.D{
  //       {Key: "name", Value: bson.D{
  //           {Key: "$regex", Value: regex},
  //       }},
  //   }},
  // }
  // log.Println("matchStage", matchStage)
  // pipeline := mongo.Pipeline{matchStage}

  // cursor, err := collection.Aggregate(context.TODO(), pipeline) 
  // if err != nil {
	// 	log.Fatal(err)
	// }
  // Decode results
  log.Println("Number of results:", cursor.RemainingBatchLength())
	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
    return
	}
  // Extract names from the results for suggestions
  // var suggestions []string
  // for _, result := range results {
  //     name := result["name"].(string)
  //     suggestions = append(suggestions, name)
  // }
  // log.Println("suggestions", suggestions)
  // Close connection
  defer database.Disconnect()

  // Return the suggestions
  c.JSON(http.StatusOK, gin.H{
    "message": "Cities suggestions",
    "searchTerm": searchTerm,
    "suggestions": results,
  })

// func handleAutocompleteOLD(c *gin.Context) {
//   //  connection
//   weatherDatabase := database.Connect()
//   collection := weatherDatabase.Collection("Cities")

//   searchTerm := c.Query("searchTerm")
//   log.Println("handleAutocomplete searchTerm", searchTerm)

//   // Fetch all documents from the collection
//   cursor, err := collection.Find(context.Background(), bson.M{})
//   if err != nil {
//     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//     return
//   }
//   defer cursor.Close(context.Background())
  
//   // Check if any documents were found
//   var suggestions []bson.M
//   for cursor.Next(context.Background()) {
//     var result bson.M
//     err := cursor.Decode(&result)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
//     name, _ := result["name"].(string)
//     if strings.HasPrefix(strings.ToLower(name), strings.ToLower(searchTerm)) {
// 			suggestions = append(suggestions, result)
// 		}
//   }

//   // Close connection
//   defer database.Disconnect()

//   // Return the suggestions
//   c.JSON(http.StatusOK, gin.H{
//     "message": "Cities suggestions",
//     "searchTerm": searchTerm,
//     "suggestions": suggestions,
//   })
}