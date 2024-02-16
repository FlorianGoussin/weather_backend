package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err error
	Client *mongo.Client
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
			return false
	}
	return true
}

// func Connect() *mongo.Client {
func Connect() {
	mongoUriEnv := func() string { 
		if isRunningInContainer() { 
			return "MONGODB_URI_CONTAINER" 
		} else {
			return "MONGODB_URI"
		}
	}()
	mongoUri := os.Getenv(mongoUriEnv)

	if mongoUri == "" {
    log.Fatal("'MONGODB_URI' environment variable not set!")
	}

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  clientOptions := options.Client().
      ApplyURI(mongoUri).
      SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil { 
		log.Fatal(err) 
	}
	Client = client
}

func Disconnect() {
	Client.Disconnect(context.Background())
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	databaseName := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}

func CollectionExists(db *mongo.Database, collectionName string) bool {
	coll, _ := db.ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: collectionName}})
	return len(coll) == 1
}