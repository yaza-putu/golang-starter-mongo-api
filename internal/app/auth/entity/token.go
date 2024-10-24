package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenType string

const ACCESS_TOKEN TokenType = "access_token"
const REFRESH_TOKEN TokenType = "refresh_token"

type Token struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type RefreshToken struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required"`
}

type UserToken struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceId  string             `json:"device_id" bson:"device_id"`
	UserId    string             `json:"user_id" bson:"user_id"`
	Token     string             `json:"token" bson:"token"`
	Device    string             `json:"device" bson:"device"`
	IP        string             `json:"ip" bson:"ip"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `json:"updated_at" bson:"updated_at"`
	TokenType TokenType
}
