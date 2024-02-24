package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Email         *string            `json:"email" validate:"email,required"`
	Password      *string            `json:"password" validate:"required"`
	Token         *string            `json:"token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_ad"`
	User_id       string             `json:"userId"`
	Refresh_token *string            `json:"refresh_token"`
	Cities        []string           `json:"cities"`
}