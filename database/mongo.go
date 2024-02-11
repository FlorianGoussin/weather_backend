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

func Connect() *mongo.Database {
	mongoUri := os.Getenv("MONGODB_URI")
  if mongoUri == "" {
    log.Fatal("'MONGODB_URI' environment variable not set!")
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  clientOptions := options.Client().
      ApplyURI(mongoUri).
      SetServerAPIOptions(serverAPIOptions)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil { 
		log.Fatal(err) 
	}
	// initializeDatabase will create Weather database and
	// the Cities collection with preloaded data
	return initializeDatabase(client)
}

func Disconnect() {
	client.Disconnect(context.Background())
}