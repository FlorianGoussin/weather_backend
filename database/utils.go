package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func GetDatabase(client *mongo.Client) *mongo.Database {
// 	databaseName := os.Getenv("DATABASE_NAME")
// 	return client.Database(databaseName)
// }

// func GetCollection(db *mongo.Database, collectionName string) *mongo.Collection {
// 	var collection *mongo.Collection = db.Collection(collectionName)
// 	return collection
// }

func CollectionExists(db *mongo.Database, collectionName string) bool {
	coll, _ := db.ListCollectionNames(context.Background(), bson.D{{Key: "name", Value: collectionName}})
	return len(coll) == 1
}

func EntryExists(filter interface{}, collection *mongo.Collection) (bool, error) {
	ctx := context.Background()

	// Perform a find operation on the collection
	err := collection.FindOne(ctx, filter).Err()

	// Check if the entry exists
	if err == mongo.ErrNoDocuments {
			return false, nil
	} else if err != nil {
			return false, err
	}

	return true, nil
}