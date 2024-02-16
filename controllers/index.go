package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	database "floriangoussin.com/weather-backend/database"
	helper "floriangoussin.com/weather-backend/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golangcompany/JWT-Authentication/models"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.GetCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
			log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
			msg = "E-Mail or Password is incorrect"
			check = false
	}
	return check, msg
}

func Register(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	var user models.User

	if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error detected while fetching the email"})
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "E-Mail already exists"})
	}

	password := HashPassword(*user.Password)
	user.Password = &password

	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Error occured while fetching the phone number"})
	}
	if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Phone Number already exists"})
	}

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
			msg := "User Details were not Saved"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
	}
	c.JSON(http.StatusOK, resultInsertionNumber)
}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
	}

	if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
	helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}
	c.JSON(http.StatusOK, foundUser)
}

// func GetUser(c *gin.Context) {
// 	userId := c.Param("user_id")
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	var user models.User
// 	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
// 	defer cancel()
// 	if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 	}
// 	c.JSON(http.StatusOK, user)
// }


func Autocomplete(c *gin.Context) {
	client := database.Client
  collection := client.Database("Weather").Collection("Cities")

  searchTerm := c.Query("searchTerm")
  log.Println("handleAutocomplete searchTerm", searchTerm)

  // TODO: Sanitize searchTerm

  // Define a filter to search for suggestions based on the searchTerm in the "name" field
  pattern := fmt.Sprintf("^%s", searchTerm)
  filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: pattern, Options: "i"}}}

  // Execute the find operation
  cursor, err := collection.Find(context.Background(), filter)
  if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
  }
  defer cursor.Close(context.Background())

  // Decode results
  log.Println("Number of results:", cursor.RemainingBatchLength())
	// var results []bson.M
	var results []models.City
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Fatal(err)
    return
	}

  // Return the suggestions
  c.JSON(http.StatusOK, gin.H{
    "message": "Cities suggestions",
    "searchTerm": searchTerm,
    "suggestions": results,
  })
}