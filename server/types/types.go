package types

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
