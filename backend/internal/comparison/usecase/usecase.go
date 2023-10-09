package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type comparisonUsecase struct {
	repo domain.ComparisonRepository
}

func NewComparisonUsecase(repo domain.ComparisonRepository) *comparisonUsecase {
	return &comparisonUsecase{repo: repo}
}

func (uc *comparisonUsecase) GetComparisons(
	ctx context.Context,
	filter *domain.ComparisonFilter,
) ([]*domain.Comparison, error) {
	comparisons, err := uc.repo.GetComparisons(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get comparisons - %w", err)
	}

	return comparisons, nil
}

func (uc *comparisonUsecase) GetComparisonById(
	ctx context.Context,
	id string,
) (*domain.Comparison, error) {
	comparison, err := uc.repo.GetComparisonById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get comparison - %w", err)
	}

	return comparison, nil
}

func (uc *comparisonUsecase) UpdateComparison(
	ctx context.Context,
	id string,
	comparison *domain.Comparison,
) error {
	comparison.Id = id
	if err := uc.repo.UpdateComparison(ctx, comparison); err != nil {
		return fmt.Errorf("failed to update comparison - %w", err)
	}

	return nil
}

func (uc *comparisonUsecase) CreateComparison(
	ctx context.Context,
	comparison *domain.Comparison,
) error {
	comparison.CreatedAt = time.Now()

	if err := uc.repo.CreateComparison(ctx, comparison); err != nil {
		return fmt.Errorf("failed to create comparison - %w", err)
	}

	return nil
}

func (uc *comparisonUsecase) DeleteComparison(ctx context.Context, id string) error {
	if err := uc.repo.DeleteComparison(ctx, id); err != nil {
		return fmt.Errorf("failed to delete comparison - %w", err)
	}

	return nil
}
