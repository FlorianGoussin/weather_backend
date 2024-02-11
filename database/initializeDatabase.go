package database

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type City struct {
	Name string `json:"name"`
	// other fields as per your JSON structure
}

// Create Weather database if there is none
func initializeDatabase(client *mongo.Client) *mongo.Database {
	ctx := context.Background()

	// Check if the database "Weather" exists
	database := client.Database("Weather")
	collections, err := database.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
			log.Fatal(err)
	}

	// Check if the "Cities" collection exists
	// If the "Weather" database or "Cities" collection doesn't exist, create them. 
	// If Weather database that doesn't exist, MongoDB will automatically create it
	var citiesCollectionExists bool
	for _, collection := range collections {
		if collection == "Cities" {
			citiesCollectionExists = true
			break
		}
	}
	if !citiesCollectionExists {
		err = database.CreateCollection(ctx, "Cities")
		if err != nil {
				log.Fatal(err)
		}
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
	data, err := os.ReadFile("cities.json")
	if err != nil {
			return err
	}
	// Parse JSON data
	var cities []City
	err = json.Unmarshal(data, &cities)
	if err != nil {
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