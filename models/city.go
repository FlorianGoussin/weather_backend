package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type City struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        *string            `json:"name" validate:"required"`
	Country     *string            `json:"country" validate:"required"`
}