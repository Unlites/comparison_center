package customoption

import (
	"context"
	"errors"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomOptionRepositoryMongo struct {
	customOptionsColl *mongo.Collection
}

type customOptionMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

func NewCustomOptionRepositoryMongo(client *mongo.Client) *CustomOptionRepositoryMongo {
	return &CustomOptionRepositoryMongo{
		customOptionsColl: client.Database("database").Collection("custom_options"),
	}
}

func (repo *CustomOptionRepositoryMongo) CreateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	_, err := repo.customOptionsColl.InsertOne(ctx, toCustomOptionMongo(customOption))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf(
				"custom option with name '%s' %w",
				customOption.Name,
				domain.ErrAlreadyExists,
			)
		}

		return fmt.Errorf("insert to mongo error: %w", err)
	}

	return nil
}

func (repo *CustomOptionRepositoryMongo) GetCustomOptionById(
	ctx context.Context,
	id string,
) (domain.CustomOption, error) {
	res := repo.customOptionsColl.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return domain.CustomOption{}, fmt.Errorf("custom option %w", domain.ErrNotFound)
		}

		return domain.CustomOption{}, fmt.Errorf("get custom option from mongo error %w", res.Err())
	}

	var com customOptionMongo
	if err := res.Decode(&com); err != nil {
		return domain.CustomOption{}, fmt.Errorf("decode mongo result error %w", err)
	}

	return toDomainCustomOption(com), nil
}

func (repo *CustomOptionRepositoryMongo) GetCustomOptions(
	ctx context.Context,
	filter domain.CustomOptionFilter,
) ([]domain.CustomOption, error) {
	opts := options.Find().
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))

	condition := bson.M{}

	if filter.Name != "" {
		condition["name"] = bson.M{"$regex": filter.Name}
	}

	cur, err := repo.customOptionsColl.Find(ctx, condition, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch custom options from mongo error: %w", err)
	}

	customOptions := make([]domain.CustomOption, 0, filter.Limit)
	for cur.Next(ctx) {
		var com customOptionMongo
		if err := cur.Decode(&com); err != nil {
			return nil, fmt.Errorf("decode mongo result error %w", err)
		}

		customOptions = append(customOptions, toDomainCustomOption(com))
	}

	return customOptions, nil
}

func (repo *CustomOptionRepositoryMongo) UpdateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	res, err := repo.customOptionsColl.UpdateOne(
		ctx,
		bson.M{"_id": customOption.Id},
		bson.M{"$set": toCustomOptionMongo(customOption)},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf(
				"custom option with name '%s' %w",
				customOption.Name,
				domain.ErrAlreadyExists,
			)
		}

		return fmt.Errorf("update at mongo error: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("custom option %w", domain.ErrNotFound)
	}

	return nil
}

func (repo *CustomOptionRepositoryMongo) DeleteCustomOption(
	ctx context.Context,
	id string,
) error {
	res, err := repo.customOptionsColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete from mongo error: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("custom option %w", domain.ErrNotFound)
	}

	return nil
}

func toCustomOptionMongo(domainCustomOption domain.CustomOption) customOptionMongo {
	return customOptionMongo{
		Id:   domainCustomOption.Id,
		Name: domainCustomOption.Name,
	}
}

func toDomainCustomOption(com customOptionMongo) domain.CustomOption {
	return domain.CustomOption{
		Id:   com.Id,
		Name: com.Name,
	}
}
