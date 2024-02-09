package models

import(
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func mongoConnect() {
	ctx := context.TODO()
  uri := "mongodb+srv://cluster0.h8m4wfl.mongodb.net/?authSource=%24external&authMechanism=MONGODB-X509&retryWrites=true&w=majority&tlsCertificateKeyFile=../X509-cert.pem"
  serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
  clientOptions := options.Client().
      ApplyURI(uri).
      SetServerAPIOptions(serverAPIOptions)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { log.Fatal(err) }
	defer client.Disconnect(ctx)
	collection := client.Database("testDB").Collection("testCol")
	docCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	fmt.Println(docCount)
}


// type City struct {
//   ID uint `json:"Id" gorm:"primary_key"`
//   Country string `json:"Country"`
//   Name string `json:"Name"`
// 	AlternativeNames pq.StringArray `gorm:"type:text[]" json:"AlternativeNames"``
// }