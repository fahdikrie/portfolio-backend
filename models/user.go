package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       			primitive.ObjectID `json:"id" xml:"id" bson:"_id,omitempty"`
	UserID	 			string  `json:"user_id"`
	Name 		 			*string `json:"name" validate:"required"`
	Fullname 			*string `json:"fullname" validate:"required"`
	Username 			*string `json:"username" validate:"required"`
	Password 			*string `json:"password,omitempty" validate:"required"`
	Token    			*string `json:"token"`
	RefreshToken  *string `json:"refresh_token"`
}
