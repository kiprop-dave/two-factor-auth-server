package models

type User struct {
	Name        string `json:"name,omitempty" validate:"required"`
	TagId       string `json:"tagId,omitempty" validate:"required"`
	PhoneNumber string `json:"phoneNumber,omitempty" validate:"required"`
	UserId      string `json:"userId,omitempty"`
}
