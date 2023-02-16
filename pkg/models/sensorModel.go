package models

type Sensor struct {
	Name   string `json:"name,omitempty" validate:"required"`
	ApiKey string `json:"apiKey" validate:"required"`
}
