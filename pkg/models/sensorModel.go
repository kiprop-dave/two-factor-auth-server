package models

type Sensor struct {
	Name   string `json:"name,omitempty" validate:"required" bson:"name"`
	ApiKey string `json:"apiKey" validate:"required" bson:"apiKey"`
}

type AuthUserRequest struct {
	Name   string `json:"name" validate:"required"`
	ApiKey string `json:"apiKey" validate:"required"`
	TagId  string `json:"tagId" validate:"required"`
}
