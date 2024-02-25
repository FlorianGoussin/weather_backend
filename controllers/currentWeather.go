package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"floriangoussin.com/weather-backend/database"
	m "floriangoussin.com/weather-backend/models"
	"github.com/gin-gonic/gin"
)

// Handlers functions
func GetAllCurrentWeatherByCity() gin.HandlerFunc { return getAllCurrentWeatherByCity }
func AddCurrentWeatherByCity() gin.HandlerFunc { return addCurrentWeatherByCity }
func RemoveCurrentWeatherByCity() gin.HandlerFunc { return removeCurrentWeatherByCity }

// WeatherResult struct to hold the result of weather request
type WeatherResult struct {
	Data  m.WeatherData
	Error error
}

var (
	weatherApiKey string
	weatherApiUrl string
)

// @Summary      Get current weather entries using user cities entries
// @Description  Get current weather entries using user cities entries
// @Tags         get all current weather
// @Accept       json
// @Produce      json
// @Param        token header string true "User JWT Token"
// @Success      200  {string} string "Successfully returned all the current weather by location entries"
// @Router       /currentWeather [get]
func getAllCurrentWeatherByCity(c *gin.Context) {
	weatherApiKey = os.Getenv("WEATHER_API_KEY")
	weatherApiUrl = os.Getenv("WEATHER_API_URL")

	// Using the JWT token find user
	token := c.Request.Header.Get("token")
	usersCollection := database.Database.Collection(database.USERS_COLLECTION)

	var user m.User
	userTokenFilter := bson.M{"token": token}
	err := usersCollection.FindOne(context.Background(), userTokenFilter).Decode(&user); 
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	// Convert user city IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, cityID := range user.Cities {
			objID, err := primitive.ObjectIDFromHex(cityID)
			if err != nil {
					log.Printf("Invalid ObjectID: %v", err)
					continue
			}
			objectIDs = append(objectIDs, objID)
	}
	if len(objectIDs) == 0 {
		log.Println("No valid ObjectIDs to search for.")
		return
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}} // find cities by IDs
	citiesCollection := database.Database.Collection(database.CITIES_COLLECTION)
	cursor, err := citiesCollection.Find(context.Background(), filter)
	if err != nil {
			log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var cities []m.City
	if err := cursor.All(context.Background(), &cities); err != nil {
			log.Fatal(err)
	}

	// Fetch current weather for all the cities
	var wg sync.WaitGroup
	wg.Add(len(cities))
	results := make(chan WeatherResult, len(cities))
	for _, city := range cities {
		fetchWeatherData(city, results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from the channel
	var weatherResults []WeatherResult
	for result := range results {
		if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
		}
		weatherResults = append(weatherResults, result)
		// Check if all results are collected
		if len(weatherResults) == len(cities) {
				break
		}
	}

	c.JSON(http.StatusOK, weatherResults)
}

// This function will be called as a go subroutine
func fetchWeatherData(city m.City, results chan<- WeatherResult, wg *sync.WaitGroup) {
	defer wg.Done()

	currentWeatherUrl := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", weatherApiUrl, weatherApiKey, *city.Name)

	// Get the current weather for the city
	response, err := http.Get(currentWeatherUrl)
	if err != nil {
			results <- WeatherResult{Error: fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)}
			return
	}
	defer response.Body.Close()

	var weatherData m.WeatherData
	if err := json.NewDecoder(response.Body).Decode(&weatherData); err != nil {
			results <- WeatherResult{Error: fmt.Errorf("failed to parse the response body: %v", err)}
			return
	}

	// Send the weather data through the channel
	results <- WeatherResult{Data: weatherData}
}

// @Summary      Add Current Weather using location and user information
// @Description  Add Current Weather using location and user information
// @Tags         add current weather location
// @Accept       json
// @Produce      json
// @Param        token header string true "User JWT Token"
// @Param        city body string true "city name"
// @Param        country body string true "country name"
// @Success      200  {string} string "Successfully returned all the current weather by location entries"
// @Router       /currentWeather [post]
func addCurrentWeatherByCity(c *gin.Context) {
	weatherApiKey = os.Getenv("WEATHER_API_KEY")
	weatherApiUrl = os.Getenv("WEATHER_API_URL")

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

	// Get the current weather for the city
	currentWeatherUrl := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", weatherApiUrl, weatherApiKey, *city.Name)
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

// @Summary      Remove Current Weather location from user
// @Description  Remove Current Weather using location and user information
// @Tags         remove current weather location
// @Accept       json
// @Produce      json
// @Param        token header string true "User JWT Token"
// @Param        city body string true "city name"
// @Param        country body string true "country name"
// @Success      200  {string} string "Successfully returned all the current weather by location entries"
// @Router       /currentWeather [delete]
func removeCurrentWeatherByCity(c *gin.Context) {
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

	update := bson.M{"$set": bson.M{"cities": updatedCityIds}}
	_, err = usersCollection.UpdateOne(context.Background(), userTokenFilter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "city deleted")
}

// Check if a user has a city id
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