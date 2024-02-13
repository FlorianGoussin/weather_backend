package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	jsoniter "github.com/json-iterator/go"
)

type City struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func dbHasCollection(db *mongo.Database, collectionName string) bool {
	coll, _ := db.ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: collectionName}})
	return len(coll) == 1
}

// Create Weather database if there is none
func Initialize(client *mongo.Client) *mongo.Database {
	database := client.Database("Weather")
	collectionName := "Cities"
	collectionExists := dbHasCollection(database, collectionName)
	if !collectionExists {
		collection := database.Collection(collectionName)

		// Insert cities from json file
		err = insertCitiesFromDataset(collection)
		if err != nil {
			log.Panic(err)
		}
		log.Println("insertCitiesFromDataset done")
		// Create index in ascending order on the name field
		if err = createIndex(collection, "name", 1); err != nil {
			log.Panic(err)
		}
		log.Println("createIndex Cities 'name' done")
	}
	// Return the Weather database
	return database
}

func insertCitiesFromDataset(collection *mongo.Collection) error {
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
	_, err = collection.InsertMany(context.Background(), cityInterfaces)
	if err != nil {
			return err
	}
	return nil // no error
}

func createIndex(collection *mongo.Collection, fieldName string, order int) error {
	indexModel := mongo.IndexModel{
    Keys: bson.D{{Key: fieldName, Value: order}},
	}
	name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
			return err
	}
	fmt.Println("Name of Index Created: " + name)
	return nil
}