package models

type User struct {
	Name        string `json:"name,omitempty" validate:"required"`
	TagId       string `json:"tagId,omitempty" validate:"required" bson:"tagId"`
	PhoneNumber string `json:"phoneNumber,omitempty" validate:"required" bson:"phoneNumber"`
	UserId      string `json:"userId,omitempty" bson:"userId"`
}
