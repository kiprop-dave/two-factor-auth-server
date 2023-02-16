package models

import (
	"time"
)

type Attempt struct {
	UserId string    `json:"userId,omitempty" validate:"required"`
	Time   time.Time `json:"time" validate:"required"`
}
