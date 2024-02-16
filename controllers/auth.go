package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	database "floriangoussin.com/weather-backend/database"
	helper "floriangoussin.com/weather-backend/helpers"
	models "floriangoussin.com/weather-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	userCollection := database.GetCollection(database.Client, "user")
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	
	var user models.User

	if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}

	validate := validator.New()
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

	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, user.User_id)
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
	userCollection := database.GetCollection(database.Client, "user")
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
	token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, foundUser.User_id)
	helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}
	c.JSON(http.StatusOK, foundUser)
}