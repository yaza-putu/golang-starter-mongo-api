package repository

import (
	"context"

	"github.com/yaza-putu/golang-starter-mongo-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

// UserInterface / **************************************************************
type User interface {
	FindByEmail(ctx context.Context, email string) (entity.User, error)
}

type userRepository struct {
	entity entity.User
}

func NewUser() *userRepository {
	return &userRepository{
		entity: entity.User{},
	}
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	e := u.entity
	err := database.Mongo.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&e)
	return e, err
}
