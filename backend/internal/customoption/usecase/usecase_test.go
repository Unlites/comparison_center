package usecase

import (
	"context"
	"testing"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/Unlites/comparison_center/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCustomOptions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		filter := domain.CustomOptionFilter{
			Limit:  2,
			Offset: 0,
		}

		returnedCustomOptions := []domain.CustomOption{
			{
				Id:   "190324fdsjfn123213",
				Name: "Speed",
			},
			{
				Id:   "303242ngpewrm40231",
				Name: "Release year",
			},
		}

		repo.On("CustomOptions", ctx, filter).Return(returnedCustomOptions, nil)

		customOptions, err := uc.CustomOptions(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, customOptions, returnedCustomOptions)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		filter := domain.CustomOptionFilter{
			Limit:  2,
			Offset: 0,
		}

		repo.On("CustomOptions", ctx, filter).Return(nil, assert.AnError)

		customOptions, err := uc.CustomOptions(ctx, filter)

		assert.Nil(t, customOptions)
		assert.Error(t, err)

		repo.AssertExpectations(t)
	})
}

func TestCreateCustomOption(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		inputCustomOption := domain.CustomOption{
			Id:   "190324fdsjfn123213",
			Name: "Speed",
		}

		repo.On("CreateCustomOption", ctx, inputCustomOption).Return(nil)
		generator.On("GenerateId").Return("190324fdsjfn123213")

		err := uc.CreateCustomOption(ctx, inputCustomOption)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		inputCustomOption := domain.CustomOption{
			Id:   "190324fdsjfn123213",
			Name: "Speed",
		}

		repo.On("CreateCustomOption", ctx, inputCustomOption).Return(assert.AnError)
		generator.On("GenerateId").Return("190324fdsjfn123213")

		err := uc.CreateCustomOption(ctx, inputCustomOption)

		assert.Error(t, err)
	})
}

func TestDeleteCustomOption(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		id := "190324fdsjfn123213"

		repo.On("DeleteCustomOption", ctx, id).Return(nil)

		err := uc.DeleteCustomOption(ctx, id)

		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		id := "190324fdsjfn123213"

		repo.On("DeleteCustomOption", ctx, id).Return(assert.AnError)

		err := uc.DeleteCustomOption(ctx, id)

		assert.Error(t, err)
	})
}

func TestCustomOptionById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		id := "190324fdsjfn123213"

		returnedCustomOption := domain.CustomOption{
			Id:   "190324fdsjfn123213",
			Name: "Speed",
		}

		repo.On("CustomOptionById", ctx, id).Return(returnedCustomOption, nil)

		customOption, err := uc.CustomOptionById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, returnedCustomOption, customOption)

		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := mocks.NewCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewCustomOptionUsecase(repo, generator)

		ctx := context.Background()

		id := "190324fdsjfn123213"

		repo.On("CustomOptionById", ctx, id).Return(domain.CustomOption{}, assert.AnError)

		customOption, err := uc.CustomOptionById(ctx, id)

		assert.Error(t, err)
		assert.Empty(t, customOption)

		repo.AssertExpectations(t)
	})
}
