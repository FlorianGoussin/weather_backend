package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"floriangoussin.com/weather-backend/database"
	m "floriangoussin.com/weather-backend/models"
	"github.com/gin-gonic/gin"
)

// Handlers functions
// func GetAllCurrentWeatherByCity() gin.HandlerFunc { return getAllCurrentWeatherByCity }
func AddCurrentWeatherByCity() gin.HandlerFunc { return addCurrentWeatherByCity }
func RemoveCurrentWeatherByCity() gin.HandlerFunc { return removeCurrentWeatherByCity }


// @Summary      Get all current weather by location entries
// @Description  Get all current weather by location entries
// @Tags         current weather location
// @Accept       json
// @Produce      json
// @Success      200  {string} string "Successfully returned all the current weather by location entries"
// @Router       /register [get]
// func getAllCurrentWeatherByCity(c *gin.Context) {
// 	token := c.GetHeader("token")

// log.Println("getAllCurrentWeatherByCity user token", token)
// 	db := mongodb.Database
// 	userCollection := db.Collection(mongodb.USERS_COLLECTION)

// 	// Use the token to find the corresponding user document
// 	filter := bson.M{"token": token}

// 	var user m.User
// 	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
// 	if err != nil {
// 			c.JSON(500, gin.H{"error": "Error finding user"})
// 			return
// 	}
// 	collectionExists := mongodb.CollectionExists(db, "userLocation")
// 	if !collectionExists {
// 		c.JSON(500, gin.H{"error": "No userLocation collection found"})
// 		return
// 	}
// 	log.Println("getAllCurrentWeatherByLocation collectionExists true!")
// 	// Perform aggregation to join User, UserLocation, and Location collections
// 	userLocationCollection := db.Collection("userLocation")
// 	pipeline := mongo.Pipeline{
// 		bson.D{{"$match", bson.D{{"location_id", user.ID}}}},
// 		bson.D{{"$lookup", bson.D{
// 				{"from", "location"},
// 				{"let", bson.D{
// 						{"location_id", "$location_id"},
// 				}},
// 				{"pipeline", bson.A{
// 						bson.D{{"$match", bson.D{
// 								{"$expr", bson.D{
// 										{"$eq", bson.A{"$location_id", "$$location_id"}},
// 								}},
// 						}}},
// 				}},
// 				{"as", "locations"},
// 		}}},
// 		bson.D{{"$unwind", "$locations"}},
// 		bson.D{{"$replaceRoot", bson.D{
// 				{"newRoot", "$locations"},
// 		}}},
// 	}
	
// 	cursor, err := userLocationCollection.Aggregate(context.Background(), pipeline)
// 	if err != nil {
// 			c.JSON(500, gin.H{"error": "Error aggregating user locations"})
// 			return
// 	}

// 	var locations []m.City
//     err = cursor.All(context.Background(), &locations)
//     if err != nil {
//         c.JSON(500, gin.H{"error": "Error decoding user locations"})
//         return
//     }

//     // At this point, 'locations' contains details of all user locations
//     c.JSON(200, locations)
// }

// @Summary      Add Current Weather using location and user information
// @Description  Add Current Weather using location and user information
// @Tags         add current weather location
// @Accept       json
// @Produce      json
// @Param        city body string true
// @Param        country body string true
// @Success      200  {string} string "Successfully returned all the current weather by location entries"
// @Router       /register [post]
func addCurrentWeatherByCity(c *gin.Context) {
	var city m.City
	if err := c.BindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Using the JWT token find user
	token := c.Request.Header.Get("token")
	usersCollection := database.Database.Collection(database.USERS_COLLECTION)

	var user m.User
	userTokenFilter := bson.M{"token": token}
	err := usersCollection.FindOne(context.Background(), userTokenFilter).Decode(&user); 
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	
	if userHasCity(&user, city.ID.Hex()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already has this city"})
		return
	}
	// Append the new city to the user's list of cities
	log.Println("addCurrentWeatherByCity city.City_id", city.ID.Hex())
	user.Cities = append(user.Cities, city.ID.Hex())

	// Update the user document in the MongoDB collection
	update := bson.M{"$set": bson.M{"cities": user.Cities}}
	_, err = usersCollection.UpdateOne(context.Background(), userTokenFilter, update)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}

	weatherApiKey := os.Getenv("WEATHER_API_KEY") // Replace this with your WeatherAPI key
	weatherApiUrl := os.Getenv("WEATHER_API_URL") // Replace this with your WeatherAPI key

	currentWeatherUrl := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", weatherApiUrl, weatherApiKey, *city.Name)

	// Get the current weather for the city
	response, err := http.Get(currentWeatherUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	// Check the status code of the response
	if response.StatusCode != http.StatusOK {
		c.JSON(response.StatusCode, gin.H{"error": "Failed to fetch data from WeatherAPI"})
		return
	}

	var data interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse the response body"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, data)
}

// @Router       /currentWeather [delete]
func removeCurrentWeatherByCity(c *gin.Context) {
	// Input: location city and country
	// Output: location id
	// Extract city data from the request or any other source
	var cityToRemove m.City
	if err := c.BindJSON(&cityToRemove); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	token := c.Request.Header.Get("token")
	usersCollection := database.Database.Collection(database.USERS_COLLECTION)
	userTokenFilter := bson.M{"token": token}
	var user m.User
	err := usersCollection.FindOne(context.Background(), userTokenFilter).Decode(&user); 
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	// Check if the user has the city to remove
	if !userHasCity(&user, cityToRemove.ID.Hex()) {
		c.JSON(http.StatusNotFound, gin.H{"error": "City not found in user's list of cities"})
		return
	}

	// Find and remove the specific city from the user's list of cities
	updatedCityIds := make([]string, 0)
	for _, userCityId := range user.Cities {
			if userCityId != cityToRemove.ID.Hex() {
					updatedCityIds = append(updatedCityIds, userCityId)
			}
	}
}

func userHasCity(user *m.User, cityId string) bool {
	if user.Cities == nil || len(user.Cities) == 0{
		return false
	}
	for _, userCityId := range user.Cities {
		if userCityId == cityId {
			return true
		}
	}
	return false
}

// func getUserByToken(c *gin.Context, token string) (*m.User, error) {
// 	userCollection := mongodb.Database.Collection(database.USERS_COLLECTION)
// 	filter := bson.M{"token": token}

// 	var user m.User
// 	err := userCollection.FindOne(context.Background(), filter).Decode(&user); 
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, m.Error{
// 				Code:    http.StatusNotFound,
// 				Message: "User not found",
// 			}
// 		}
// 		return nil, m.Error{
// 			Code:    http.StatusInternalServerError,
// 			Message: err.Error(), // Pass the original error message
// 		}
// 	}
// 	return &user, nil
// }