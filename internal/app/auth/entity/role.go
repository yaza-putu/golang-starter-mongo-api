package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ADM = "SF%YHm8-XJ^}"
const USR = "Ts6W0l2EU8&v"

type Role struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt primitive.DateTime `json:"-" bson:"created_at"`
	UpdatedAt primitive.DateTime `json:"-" bson:"updated_at"`
}

type Roles []Role
