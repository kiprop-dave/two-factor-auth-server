package types

import (
	jwt "github.com/golang-jwt/jwt/v5"
)

type RegistrationRequest struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	TagId    string `json:"tagId" bson:"tagId"`
}

type RegistrationResponse struct {
	TwoFaQrUri string `json:"twoFaQrUri,omitempty" bson:"twoFaQrUri"`
	ID         string `json:"id,omitempty" bson:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserClaims struct {
	Email string `json:"email"`
}

type AdminClaims struct {
	User                 UserClaims `json:"user"`
	jwt.RegisteredClaims `json:"claims"`
}

type RfidCheckRequest struct {
	TagId  string `json:"tagId"`
	ApiKey string `json:"apiKey"`
}

type RfidCheckResponse struct {
	EntryAttemptId string `json:"entryAttemptId"`
	Role           string `json:"role"`
}

type TwoFaRequest struct {
	EntryAttemptId string `json:"entryAttemptId"`
	ApiKey         string `json:"apiKey"`
	TOTP           string `json:"totp"`
}

type TwoFaResponse struct {
	Success bool `json:"success"`
}
