package repository

import (
	"context"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/mock"
)

type objectRepositoryMock struct {
	mock.Mock
}

func NewObjectRepositoryMock() *objectRepositoryMock {
	return &objectRepositoryMock{}
}

func (repo *objectRepositoryMock) GetObjects(
	ctx context.Context,
	filter domain.ObjectFilter,
) ([]domain.Object, error) {
	args := repo.Called(ctx, filter)

	ret, err := args.Get(0), args.Error(1)

	var objects []domain.Object

	if ret != nil {
		objects = ret.([]domain.Object)
	}

	return objects, err
}

func (repo *objectRepositoryMock) GetObjectById(ctx context.Context, id string) (domain.Object, error) {
	args := repo.Called(ctx, id)

	ret, err := args.Get(0), args.Error(1)

	var object domain.Object

	if ret != nil {
		object = ret.(domain.Object)
	}

	return object, err
}

func (repo *objectRepositoryMock) UpdateObject(ctx context.Context, object domain.Object) error {
	args := repo.Called(ctx, object)

	return args.Error(0)
}

func (repo *objectRepositoryMock) CreateObject(ctx context.Context, object domain.Object) error {
	args := repo.Called(ctx, object)

	return args.Error(0)
}

func (repo *objectRepositoryMock) DeleteObject(ctx context.Context, id string) error {
	args := repo.Called(ctx, id)

	return args.Error(0)
}
