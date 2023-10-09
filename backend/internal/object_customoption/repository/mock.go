package repository

import (
	"context"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type objectCustomOptionRepositoryMock struct {
	mock.Mock
}

func NewObjectCustomOptionRepositoryMock() *objectCustomOptionRepositoryMock {
	return &objectCustomOptionRepositoryMock{}
}

func (repo *objectCustomOptionRepositoryMock) GetObjectCustomOptionsByObjectId(
	ctx context.Context,
	objectId string,
) ([]*domain.ObjectCustomOption, error) {
	args := repo.Called(ctx, objectId)

	ret, err := args.Get(0), args.Error(1)

	var objectCustomOptions []*domain.ObjectCustomOption

	if ret != nil {
		objectCustomOptions = ret.([]*domain.ObjectCustomOption)
	}

	return objectCustomOptions, err
}

func (repo *objectCustomOptionRepositoryMock) AddObjectCustomOption(
	ctx context.Context,
	objectCustomOption *domain.ObjectCustomOption,
) error {
	args := repo.Called(ctx, objectCustomOption)

	return args.Error(0)
}

func (repo *objectCustomOptionRepositoryMock) UpdateObjectCustomOption(
	ctx context.Context,
	objectCustomOption *domain.ObjectCustomOption,
) error {
	args := repo.Called(ctx, objectCustomOption)

	return args.Error(0)
}
