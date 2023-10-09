package usecase

import (
	"context"
	"testing"

	"github.com/Unlites/comparison_center/backend/internal/customoption/repository"
	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetCustomOptions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		returnedCustomOptions := []*domain.CustomOption{
			{
				Id:   "190324fdsjfn123213",
				Name: "Speed",
			},
			{
				Id:   "303242ngpewrm40231",
				Name: "Release year",
			},
		}

		repo := repository.NewCustomOptionRepositoryMock()
		uc := NewCustomOptionUsecase(repo)

		ctx := context.Background()
		filter := &domain.CustomOptionFilter{
			Limit:  2,
			Offset: 0,
		}

		repo.On("GetCustomOptions", ctx, filter).Return(returnedCustomOptions, nil)

		customOptions, err := uc.GetCustomOptions(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, customOptions, returnedCustomOptions)
		repo.AssertExpectations(t)
	})

}
