package repository

import (
	"context"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type comparisonRepositoryMock struct {
	mock.Mock
}

func NewComparisonRepositoryMock() *comparisonRepositoryMock {
	return &comparisonRepositoryMock{}
}

func (repo *comparisonRepositoryMock) GetComparisons(
	ctx context.Context,
	filter domain.ComparisonFilter,
) ([]domain.Comparison, error) {
	args := repo.Called(ctx, filter)

	ret, err := args.Get(0), args.Error(1)

	var comparisons []domain.Comparison

	if ret != nil {
		comparisons = ret.([]domain.Comparison)
	}

	return comparisons, err
}

func (repo *comparisonRepositoryMock) GetComparisonById(
	ctx context.Context,
	id string,
) (domain.Comparison, error) {
	args := repo.Called(ctx, id)

	ret, err := args.Get(0), args.Error(1)

	var comparison domain.Comparison

	if ret != nil {
		comparison = ret.(domain.Comparison)
	}

	return comparison, err
}

func (repo *comparisonRepositoryMock) UpdateComparison(
	ctx context.Context,
	comparison domain.Comparison,
) error {
	args := repo.Called(ctx, comparison)

	return args.Error(0)
}

func (repo *comparisonRepositoryMock) CreateComparison(
	ctx context.Context,
	comparison domain.Comparison,
) error {
	args := repo.Called(ctx, comparison)

	return args.Error(0)
}

func (repo *comparisonRepositoryMock) DeleteComparison(
	ctx context.Context,
	id string,
) error {
	args := repo.Called(ctx, id)

	return args.Error(0)
}
