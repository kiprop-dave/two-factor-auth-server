package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Attempt struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	UserId string             `json:"userId,omitempty" validate:"required"`
	Time   time.Time          `json:"time" validate:"required"`
}
