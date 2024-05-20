package repository

import (
	"context"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ObjectCustomOptionRepositoryMongo struct {
	objectCustomOptionsColl *mongo.Collection
}

func NewObjectCustomOptionRepositoryMongo(client *mongo.Client) *ObjectCustomOptionRepositoryMongo {
	return &ObjectCustomOptionRepositoryMongo{
		objectCustomOptionsColl: client.Database("database").Collection("object_custom_options"),
	}
}

type objectCustomOptionMongo struct {
	ObjectId       string `bson:"object_id"`
	CustomOptionId string `bson:"custom_option_id"`
	Value          string `bson:"value"`
}

func (repo *ObjectCustomOptionRepositoryMongo) GetObjectCustomOptionsByObjectId(
	ctx context.Context,
	objectId string,
) ([]domain.ObjectCustomOption, error) {
	cur, err := repo.objectCustomOptionsColl.Find(ctx, bson.M{"object_id": objectId})
	if err != nil {
		return nil, fmt.Errorf("fetch object custom options from mongo error: %w", err)
	}

	objCustomOptions := make([]domain.ObjectCustomOption, 0)
	for cur.Next(ctx) {
		var ocom objectCustomOptionMongo
		if err := cur.Decode(ocom); err != nil {
			return nil, fmt.Errorf("decode mongo result error %w", err)
		}

		objCustomOptions = append(objCustomOptions, toDomainObjectCustomOption(ocom))
	}

	return objCustomOptions, nil
}

func (repo *ObjectCustomOptionRepositoryMongo) AddObjectCustomOption(
	ctx context.Context,
	objectCustomOption domain.ObjectCustomOption,
) error {
	_, err := repo.objectCustomOptionsColl.InsertOne(
		ctx,
		toObjectCustomOptionMongo(objectCustomOption),
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf(
				"option with object id '%s' and custom option id '%s' %w",
				objectCustomOption.ObjectId,
				objectCustomOption.CustomOptionId,
				domain.ErrAlreadyExists,
			)
		}
		return fmt.Errorf("insert to mongo error: %w", err)
	}

	return nil
}

func (repo *ObjectCustomOptionRepositoryMongo) UpdateObjectCustomOption(
	ctx context.Context,
	objectCustomOption domain.ObjectCustomOption,
) error {
	res, err := repo.objectCustomOptionsColl.UpdateOne(
		ctx,
		bson.M{
			"object_id":        objectCustomOption.ObjectId,
			"custom_option_id": objectCustomOption.CustomOptionId,
		},
		bson.M{"$set": toObjectCustomOptionMongo(objectCustomOption)},
	)
	if err != nil {
		return fmt.Errorf("update at mongo error: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("object custom option %w", domain.ErrNotFound)
	}

	return nil
}

func toDomainObjectCustomOption(ocom objectCustomOptionMongo) domain.ObjectCustomOption {
	return domain.ObjectCustomOption{
		ObjectId:       ocom.ObjectId,
		CustomOptionId: ocom.CustomOptionId,
		Value:          ocom.Value,
	}
}

func toObjectCustomOptionMongo(objectCustomOption domain.ObjectCustomOption) objectCustomOptionMongo {
	return objectCustomOptionMongo{
		ObjectId:       objectCustomOption.ObjectId,
		CustomOptionId: objectCustomOption.CustomOptionId,
		Value:          objectCustomOption.Value,
	}
}
