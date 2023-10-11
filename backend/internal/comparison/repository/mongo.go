package repository

import (
	"context"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type comparisonRepositoryMongo struct {
	client *mongo.Client
}

func NewComparisonRepositoryMongo(client *mongo.Client) *comparisonRepositoryMongo {
	return &comparisonRepositoryMongo{client: client}
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
	return fmt.Errorf("not implemented")
}

func (repo *comparisonRepositoryMongo) DeleteComparison(
	ctx context.Context,
	id string,
) error {
	return fmt.Errorf("not implemented")
}
