package database

import(
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
    log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  clientOptions := options.Client().
      ApplyURI(mongoUri).
      SetServerAPIOptions(serverAPIOptions)
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil { 
		log.Fatal(err) 
	}
	return client.Database("Weather")
}

func Disconnect() {
	client.Disconnect(context.Background())
}