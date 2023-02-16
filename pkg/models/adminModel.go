package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
}
