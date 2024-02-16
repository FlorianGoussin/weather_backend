package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	First_name *string `json:"first_name" validate:"required,min:2,max:100"`
	Last_name *string `json:"last_name" validate:"required,min:2,max:100"`
	Email *string `json:"email" validate:"email,required"`
	Password *string `json:"Password" validate:"required,min:6"`
	Token *string `json:"token"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_ad"`
	User_id string `json:"userId"`
	Refresh_token *string `json:"refresh_token"`
}