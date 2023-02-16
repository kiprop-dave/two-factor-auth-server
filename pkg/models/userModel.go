package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name,omitempty" validate:"required"`
	TagId       string             `json:"tagId,omitempty" validate:"required"`
	PhoneNumber string             `json:"phoneNumber,omitempty" validate:"required"`
	UserId      string             `json:"userId,omitempty" validate:"required"`
}
