package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	database "floriangoussin.com/weather-backend/database"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email string
	First_name string
	Last_name string
	Uid string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.GetCollection(database.Client, "user")
var SECRET_KEY string = os.Getenv("JWT_SECRET_KEY")

// Generates a token and a refresh token
func GenerateAllTokens(email string, uid string) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	signedRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return signedToken, signedRefreshToken, err
}

// Validates the token. 
// Returns the claims if token is valid
// Returns error message if token is not valid
func ValidateToken(signedToken string)(claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
    msg := fmt.Sprintf("failed to parse token: %s", err.Error())
		log.Println(msg)
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg := fmt.Sprintf("the token is invalid: %s", err.Error())
		log.Println(msg)
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg := fmt.Sprintf("Your Token is not Valid anymore: %s", err.Error())
		log.Println(msg)
		return nil, msg
	}
	fmt.Println(claims, msg)
	return claims, msg
}

// Updates the user's token in the database
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}