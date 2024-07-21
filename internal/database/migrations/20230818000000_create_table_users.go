package migrations

import (
	"context"

	"github.com/yaza-putu/golang-starter-mongo-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/database"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/pkg/encrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	const collectionName = "users"
	database.MigrationRegister(func(context context.Context, db *mongo.Database) error {
		user := entity.User{
			ID:       primitive.NewObjectID(),
			Name:     "admin",
			Email:    "admin@mail.com",
			Password: encrypt.Bcrypt("Password1"),
		}
		_, err := db.Collection(collectionName).InsertOne(context, user)
		return err
	}, func(context context.Context, db *mongo.Database) error { // drop collection
		return db.Collection(collectionName).Drop(context)
	})
}
