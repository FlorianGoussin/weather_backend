package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	err error
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
			return false
	}
	return true
}

func Connect() *mongo.Database {
	mongoUriEnv := func() string { 
		if isRunningInContainer() { 
			return "MONGODB_URI_CONTAINER" 
		} else {
			return "MONGODB_URI"
		}
	}()
	mongoUri := os.Getenv(mongoUriEnv)
	
	mongoUsername := os.Getenv("MONGO_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_ROOT_PASSWORD")
	if mongoUri == "" {
    log.Fatal("'MONGODB_URI' environment variable not set!")
	} else if mongoUsername == "" {
    log.Fatal("'MONGO_ROOT_USERNAME' environment variable not set!")
	} else if mongoPassword == "" {
    log.Fatal("'MONGO_ROOT_PASSWORD' environment variable not set!")
	}

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  clientOptions := options.Client().
      ApplyURI(mongoUri).
			// SetAuth(options.Credential{
			// 	Username: mongoUsername,
			// 	Password: mongoPassword,
			// }).
      SetServerAPIOptions(serverAPIOptions)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil { 
		log.Fatal(err) 
	}
	// initializeDatabase will create Weather database and
	// the Cities collection with preloaded data
	return initializeDatabase(client)
	// return client.Database("Weather")
}

func Disconnect() {
	client.Disconnect(context.Background())
}