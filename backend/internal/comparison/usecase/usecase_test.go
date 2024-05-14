package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/comparison/repository"
	"github.com/Unlites/comparison_center/backend/internal/domain"
	g "github.com/Unlites/comparison_center/backend/pkg/generator"
	"github.com/stretchr/testify/assert"
)

func TestGetComparisons(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()
		returnedComparisons := []domain.Comparison{
			{
				Id:              "85434230werhuhi123912304",
				Name:            "Cars",
				CreatedAt:       time.Now(),
				CustomOptionIds: []string{"43294320fdsfnj13213", "3240312rnwjnj49329"},
			},
			{
				Id:              "32492349mkfdsmfks234",
				Name:            "Computers",
				CreatedAt:       time.Now(),
				CustomOptionIds: []string{"012332432dfsdsjof21312321", "91230123021sdfdsf13123"},
			},
		}

		filter := domain.ComparisonFilter{
			Limit:  2,
			Offset: 0,
		}

		repo.On("GetComparisons", ctx, filter).Return(returnedComparisons, nil)

		comparisons, err := uc.GetComparisons(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, comparisons, returnedComparisons)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()
		filter := domain.ComparisonFilter{
			Limit:  2,
			Offset: 0,
		}

		repo.On("GetComparisons", ctx, filter).Return(nil, errors.New("some error"))

		comparisons, err := uc.GetComparisons(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, comparisons)
		repo.AssertExpectations(t)
	})
}

func TestGetComparisonById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		returnedComparison := domain.Comparison{
			Id:              "85434230werhuhi123912304",
			Name:            "Cars",
			CreatedAt:       time.Now(),
			CustomOptionIds: []string{"43294320fdsfnj13213", "3240312rnwjnj49329"},
		}

		ctx := context.Background()
		id := "85434230werhuhi123912304"

		repo.On("GetComparisonById", ctx, id).Return(returnedComparison, nil)

		comparison, err := uc.GetComparisonById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, comparison, returnedComparison)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()
		id := "213213ewrwe9423432"

		repo.On("GetComparisonById", ctx, id).Return(nil, domain.ErrNotFound)

		comparison, err := uc.GetComparisonById(ctx, id)

		assert.Error(t, err)
		assert.Empty(t, comparison)
		repo.AssertExpectations(t)
	})
}

func TestCreateComparison(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		inputComparison := domain.Comparison{
			Name:            "Cars",
			CustomOptionIds: []string{"3332415fdsfsd31231", "5412asdsa131231`"},
		}

		changedComparison := domain.Comparison{
			Id:              "49234991asdsanjd12305",
			Name:            inputComparison.Name,
			CustomOptionIds: inputComparison.CustomOptionIds,
			CreatedAt:       time.Now(),
		}

		ctx := context.Background()

		generator.On("GenerateId").Return("49234991asdsanjd12305")
		repo.On("CreateComparison", ctx, changedComparison).Return(nil)

		err := uc.CreateComparison(ctx, inputComparison)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		inputComparison := domain.Comparison{
			Name:            "Cars",
			CustomOptionIds: []string{"432432sadas5433da", "349fsda32bfsd21d"},
		}

		changedComparison := domain.Comparison{
			Id:              "32939fsdfsdf912312",
			Name:            inputComparison.Name,
			CustomOptionIds: inputComparison.CustomOptionIds,
			CreatedAt:       time.Now(),
		}

		ctx := context.Background()

		generator.On("GenerateId").Return("32939fsdfsdf912312")
		repo.On("CreateComparison", ctx, changedComparison).Return(errors.New("some error"))

		err := uc.CreateComparison(ctx, inputComparison)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdateComparison(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()

		inputComparison := domain.Comparison{
			Name:            "Cars",
			CustomOptionIds: []string{"23491239dqwe14sddsf", "74329fdsfsdwe13123q"},
		}

		id := "94232dsadas21313ddsa"

		returnedComparison := domain.Comparison{
			Id:              id,
			Name:            "Cars",
			CreatedAt:       time.Now(),
			CustomOptionIds: []string{"43294320fdsfnj13213", "3240312rnwjnj49329"},
		}

		changedComparison := domain.Comparison{
			Id:              id,
			Name:            inputComparison.Name,
			CreatedAt:       returnedComparison.CreatedAt,
			CustomOptionIds: inputComparison.CustomOptionIds,
		}

		repo.On("GetComparisonById", ctx, id).Return(returnedComparison, nil)
		repo.On("UpdateComparison", ctx, changedComparison).Return(nil)

		err := uc.UpdateComparison(ctx, id, inputComparison)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()

		inputComparison := domain.Comparison{
			Name:            "Cars",
			CustomOptionIds: []string{"23491239dqwe14sddsf", "74329fdsfsdwe13123q"},
		}

		id := "48213asd9332ewqse328"

		returnedComparison := domain.Comparison{
			Id:              id,
			Name:            "Cars",
			CreatedAt:       time.Now(),
			CustomOptionIds: []string{"43294320fdsfnj13213", "3240312rnwjnj49329"},
		}

		changedComparison := domain.Comparison{
			Id:              id,
			Name:            inputComparison.Name,
			CreatedAt:       returnedComparison.CreatedAt,
			CustomOptionIds: inputComparison.CustomOptionIds,
		}

		repo.On("GetComparisonById", ctx, id).Return(returnedComparison, nil)
		repo.On("UpdateComparison", ctx, changedComparison).Return(errors.New("some error"))

		err := uc.UpdateComparison(ctx, id, inputComparison)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestDeleteComparison(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()
		id := "34543dfsdfj32432jewr"

		repo.On("DeleteComparison", ctx, id).Return(nil)

		err := uc.DeleteComparison(ctx, id)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		repo := repository.NewComparisonRepositoryMock()
		generator := g.NewMockGenerator()
		uc := NewComparisonUsecase(repo, generator)

		ctx := context.Background()
		id := "92133easd123srewr132"

		repo.On("DeleteComparison", ctx, id).Return(errors.New("some error"))

		err := uc.DeleteComparison(ctx, id)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
