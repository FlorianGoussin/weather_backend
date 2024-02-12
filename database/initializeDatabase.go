package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	jsoniter "github.com/json-iterator/go"
)

type City struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

// Create Weather database if there is none
func initializeDatabase(client *mongo.Client) *mongo.Database {
	database := client.Database("Weather")
	collection := database.Collection("Cities") // Create collection if !exists
	count, _ := collection.CountDocuments(context.Background(), nil)
	if (count == 0) {
		// Insert cities from json file
		err = insertCitiesFromDataset(client)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Create and add data to collection 'Cities'")
	}
	// Return the Weather database
	return database
}

func insertCitiesFromDataset(client *mongo.Client) error {
	// Read cities.json file
	data, err := os.ReadFile("data.json")
	if err != nil {
			return err
	}
	var cities []City
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(data, &cities); err != nil {
		return err
	}
	
	// Convert slice of City to slice of interface{}
	var cityInterfaces []interface{}
	for _, city := range cities {
			cityInterfaces = append(cityInterfaces, city)
	}
	// Insert cities into MongoDB
	collection := client.Database("Weather").Collection("Cities")
	_, err = collection.InsertMany(context.Background(), cityInterfaces)
	if err != nil {
			return err
	}
	return nil // no error
}