package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err error
)

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
			return false
	}
	return true
}

func Connect() *mongo.Client {
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
	return client
}

func Disconnect(client *mongo.Client) {
	client.Disconnect(context.Background())
}