package comparison

import (
	"context"
	"fmt"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type ComparisonUsecase struct {
	repo        ComparisonRepository
	idGenerator IdGenerator
}

type ComparisonRepository interface {
	GetComparisons(ctx context.Context, filter domain.ComparisonFilter) ([]domain.Comparison, error)
	GetComparisonById(ctx context.Context, id string) (domain.Comparison, error)
	UpdateComparison(ctx context.Context, comparison domain.Comparison) error
	CreateComparison(ctx context.Context, comparison domain.Comparison) error
	DeleteComparison(ctx context.Context, id string) error
}

type IdGenerator interface {
	GenerateId() string
}

func NewComparisonUsecase(
	repo ComparisonRepository,
	idGenerator IdGenerator,
) *ComparisonUsecase {
	return &ComparisonUsecase{
		repo:        repo,
		idGenerator: idGenerator,
	}
}

func (uc *ComparisonUsecase) GetComparisons(
	ctx context.Context,
	filter domain.ComparisonFilter,
) ([]domain.Comparison, error) {
	comparisons, err := uc.repo.GetComparisons(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get comparisons - %w", err)
	}

	return comparisons, nil
}

func (uc *ComparisonUsecase) GetComparisonById(
	ctx context.Context,
	id string,
) (domain.Comparison, error) {
	comparison, err := uc.repo.GetComparisonById(ctx, id)
	if err != nil {
		return domain.Comparison{}, fmt.Errorf("failed to get comparison - %w", err)
	}

	return comparison, nil
}

func (uc *ComparisonUsecase) UpdateComparison(
	ctx context.Context,
	id string,
	comparison domain.Comparison,
) error {
	existingComparison, err := uc.repo.GetComparisonById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing comparison - %w", err)
	}

	comparison.Id = existingComparison.Id
	comparison.CreatedAt = existingComparison.CreatedAt

	if err := uc.repo.UpdateComparison(ctx, comparison); err != nil {
		return fmt.Errorf("failed to update comparison - %w", err)
	}

	return nil
}

func (uc *ComparisonUsecase) CreateComparison(
	ctx context.Context,
	comparison domain.Comparison,
) error {
	comparison.Id = uc.idGenerator.GenerateId()
	comparison.CreatedAt = time.Now()

	if err := uc.repo.CreateComparison(ctx, comparison); err != nil {
		return fmt.Errorf("failed to create comparison - %w", err)
	}

	return nil
}

func (uc *ComparisonUsecase) DeleteComparison(ctx context.Context, id string) error {
	if err := uc.repo.DeleteComparison(ctx, id); err != nil {
		return fmt.Errorf("failed to delete comparison - %w", err)
	}

	return nil
}
