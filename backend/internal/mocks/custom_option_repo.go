package mocks

import (
	"context"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CustomOptionRepositoryMock struct {
	mock.Mock
}

func NewCustomOptionRepositoryMock() *CustomOptionRepositoryMock {
	return &CustomOptionRepositoryMock{}
}

func (repo *CustomOptionRepositoryMock) CustomOptions(
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

func (repo *CustomOptionRepositoryMock) CustomOptionById(
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

func (repo *CustomOptionRepositoryMock) UpdateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	args := repo.Called(ctx, customOption)

	return args.Error(0)
}

func (repo *CustomOptionRepositoryMock) CreateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	args := repo.Called(ctx, customOption)

	return args.Error(0)
}

func (repo *CustomOptionRepositoryMock) DeleteCustomOption(
	ctx context.Context,
	id string,
) error {
	args := repo.Called(ctx, id)

	return args.Error(0)
}
