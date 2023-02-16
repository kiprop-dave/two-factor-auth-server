package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Sensor struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	Name   string             `json:"name,omitempty" validate:"required"`
	ApiKey string             `json:"apiKey" validate:"required"`
}
