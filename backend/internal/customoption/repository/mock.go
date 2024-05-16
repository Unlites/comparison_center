package repository

import (
	"context"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type customOptionRepositoryMock struct {
	mock.Mock
}

func NewCustomOptionRepositoryMock() *customOptionRepositoryMock {
	return &customOptionRepositoryMock{}
}

func (repo *customOptionRepositoryMock) GetCustomOptions(
	ctx context.Context,
	filter domain.CustomOptionFilter,
) ([]domain.CustomOption, error) {
	args := repo.Called(ctx, filter)

	ret, err := args.Get(0), args.Error(1)

	var customOptions []domain.CustomOption

	if ret != nil {
		customOptions = ret.([]domain.CustomOption)
	}

	return customOptions, err
}

func (repo *customOptionRepositoryMock) GetCustomOptionById(
	ctx context.Context,
	id string,
) (domain.CustomOption, error) {
	args := repo.Called(ctx, id)

	ret, err := args.Get(0), args.Error(1)

	var customOption domain.CustomOption

	if ret != nil {
		customOption = ret.(domain.CustomOption)
	}

	return customOption, err
}

func (repo *customOptionRepositoryMock) UpdateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	args := repo.Called(ctx, customOption)

	return args.Error(0)
}

func (repo *customOptionRepositoryMock) CreateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	args := repo.Called(ctx, customOption)

	return args.Error(0)
}

func (repo *customOptionRepositoryMock) DeleteCustomOption(
	ctx context.Context,
	id string,
) error {
	args := repo.Called(ctx, id)

	return args.Error(0)
}
