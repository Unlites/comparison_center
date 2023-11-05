package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type comparisonRepositoryMongo struct {
	comparisonsColl *mongo.Collection
}

type comparisonMongo struct {
	Id              string    `bson:"_id"`
	Name            string    `bson:"name"`
	CreatedAt       time.Time `bson:"created_at"`
	CustomOptionIds []string  `bson:"custom_option_ids"`
}

func NewComparisonRepositoryMongo(client *mongo.Client) *comparisonRepositoryMongo {
	return &comparisonRepositoryMongo{
		comparisonsColl: client.Database("database").Collection("comparisons"),
	}
}

func (repo *comparisonRepositoryMongo) GetComparisons(
	ctx context.Context,
	filter *domain.ComparisonFilter,
) ([]*domain.Comparison, error) {
	opts := options.Find().
		SetSort(bson.M{filter.OrderBy: 1}).
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))

	cur, err := repo.comparisonsColl.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch comparisons from mongo error: %w", err)
	}

	comparisons := make([]*domain.Comparison, 0, filter.Limit)
	for cur.Next(ctx) {
		cm := new(comparisonMongo)
		if err := cur.Decode(cm); err != nil {
			return nil, fmt.Errorf("decode mongo result error %w", err)
		}

		comparisons = append(comparisons, toDomainComparison(cm))
	}

	return comparisons, nil
}

func (repo *comparisonRepositoryMongo) GetComparisonById(
	ctx context.Context,
	id string,
) (*domain.Comparison, error) {
	res := repo.comparisonsColl.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("comparison %w", domain.ErrNotFound)
		}

		return nil, fmt.Errorf("get comparison from mongo error %w", res.Err())
	}

	cm := new(comparisonMongo)
	if err := res.Decode(cm); err != nil {
		return nil, fmt.Errorf("decode mongo result error %w", err)
	}

	return toDomainComparison(cm), nil
}

func (repo *comparisonRepositoryMongo) UpdateComparison(
	ctx context.Context,
	comparison *domain.Comparison,
) error {
	res, err := repo.comparisonsColl.UpdateOne(
		ctx,
		bson.M{"_id": comparison.Id},
		bson.M{"$set": toComparisonMongo(comparison)},
	)
	if err != nil {
		return fmt.Errorf("update at mongo error: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("comparison %w", domain.ErrNotFound)
	}

	return nil
}

func (repo *comparisonRepositoryMongo) CreateComparison(
	ctx context.Context,
	comparison *domain.Comparison,
) error {
	err := repo.comparisonsColl.FindOne(ctx, bson.M{
		"name": comparison.Name,
	}).Err()

	if err == nil {
		return fmt.Errorf(
			"comparison '%s' %w",
			comparison.Name,
			domain.ErrAlreadyExists,
		)
	} else {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("check existence of comparison in mongo error: %w", err)
		}
	}

	_, err = repo.comparisonsColl.InsertOne(ctx, toComparisonMongo(comparison))
	if err != nil {
		return fmt.Errorf("insert to mongo error: %w", err)
	}

	return nil
}

func (repo *comparisonRepositoryMongo) DeleteComparison(
	ctx context.Context,
	id string,
) error {
	res, err := repo.comparisonsColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete from mongo error: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("comparison %w", domain.ErrNotFound)
	}

	return nil
}

func toComparisonMongo(domainComparison *domain.Comparison) *comparisonMongo {
	slog.Info(fmt.Sprintf("%v", domainComparison))
	return &comparisonMongo{
		Id:              domainComparison.Id,
		Name:            domainComparison.Name,
		CreatedAt:       domainComparison.CreatedAt,
		CustomOptionIds: domainComparison.CustomOptionIds,
	}
}

func toDomainComparison(comparisonMongo *comparisonMongo) *domain.Comparison {
	return &domain.Comparison{
		Id:              comparisonMongo.Id,
		Name:            comparisonMongo.Name,
		CreatedAt:       comparisonMongo.CreatedAt,
		CustomOptionIds: comparisonMongo.CustomOptionIds,
	}
}
