package repository

import (
	"context"
	"time"

	"github.com/yaza-putu/golang-starter-mongo-api/internal/app/auth/entity"
	"github.com/yaza-putu/golang-starter-mongo-api/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserToken interface {
		Create(ctx context.Context, data entity.UserToken) (entity.UserToken, error)
		FindByDeviceId(ctx context.Context, deviceId string) (entity.UserToken, error)
		Revoke(ctx context.Context, eviceId string) error
		Update(ctx context.Context, deviceId string, data entity.UserToken) (entity.UserToken, error)
	}
	userTokenRepository struct {
		entity     entity.UserToken
		collection string
	}
)

func NewUserToken() *userTokenRepository {
	return &userTokenRepository{
		entity:     entity.UserToken{},
		collection: "user_tokens",
	}
}

func (u *userTokenRepository) Create(ctx context.Context, data entity.UserToken) (entity.UserToken, error) {
	// find by userid
	e := u.entity
	database.Mongo.Collection(u.collection).FindOne(ctx, bson.M{"user_id": data.UserId, "ip": data.IP, "device": data.Device}).Decode(&e)
	// update
	if !e.ID.IsZero() {
		data.ID = primitive.NilObjectID

		_, err := database.Mongo.Collection(u.collection).UpdateOne(ctx, bson.M{"_id": e.ID}, bson.M{
			"$set": data,
		})

		return data, err
	}
	// create
	data.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	data.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := database.Mongo.Collection(u.collection).InsertOne(ctx, &data)
	return data, err
}

func (u *userTokenRepository) Revoke(ctx context.Context, deviceId string) error {
	_, err := database.Mongo.Collection(u.collection).DeleteOne(ctx, bson.M{"device_id": deviceId})
	return err
}

func (u *userTokenRepository) FindByDeviceId(ctx context.Context, deviceId string) (entity.UserToken, error) {
	e := u.entity
	err := database.Mongo.Collection(u.collection).FindOne(ctx, bson.M{"device_id": deviceId}).Decode(&e)
	return e, err
}

func (u *userTokenRepository) Update(ctx context.Context, id primitive.ObjectID, data entity.UserToken) (entity.UserToken, error) {
	data.ID = primitive.NilObjectID
	data.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := database.Mongo.Collection(u.collection).UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": data,
	})
	return data, err
}
