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
)

var errMongoNotFound = errors.New("no documents in result")

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
	return nil, fmt.Errorf("not implemented")
}

func (repo *comparisonRepositoryMongo) GetComparisonById(
	ctx context.Context,
	id string,
) (*domain.Comparison, error) {
	return nil, fmt.Errorf("not implemented")
}

func (repo *comparisonRepositoryMongo) UpdateComparison(
	ctx context.Context,
	comparison *domain.Comparison,
) error {
	return fmt.Errorf("not implemented")
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
	return fmt.Errorf("not implemented")
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
