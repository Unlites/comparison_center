package usecase

import (
	"context"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type customOptionUsecase struct {
	repo domain.CustomOptionRepository
}

func NewCustomOptionUsecase(repo domain.CustomOptionRepository) *customOptionUsecase {
	return &customOptionUsecase{repo: repo}
}

func (uc *customOptionUsecase) GetCustomOptions(
	ctx context.Context,
	filter *domain.CustomOptionFilter,
) ([]*domain.CustomOption, error) {
	customOptions, err := uc.repo.GetCustomOptions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom options - %w", err)
	}

	return customOptions, nil
}

func (uc *customOptionUsecase) GetCustomOptionById(
	ctx context.Context,
	id string,
) (*domain.CustomOption, error) {
	customOption, err := uc.repo.GetCustomOptionById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom option - %w", err)
	}

	return customOption, nil
}

func (uc *customOptionUsecase) UpdateCustomOption(
	ctx context.Context,
	id string,
	customOption *domain.CustomOption,
) error {
	existingCustomOption, err := uc.repo.GetCustomOptionById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing custom option - %w", err)
	}

	customOption.Id = existingCustomOption.Id

	if err := uc.repo.UpdateCustomOption(ctx, customOption); err != nil {
		return fmt.Errorf("failed to update custom option - %w", err)
	}

	return nil
}

func (uc *customOptionUsecase) CreateCustomOption(
	ctx context.Context,
	customOption *domain.CustomOption,
) error {
	if err := uc.repo.CreateCustomOption(ctx, customOption); err != nil {
		return fmt.Errorf("failed to create custom option - %w", err)
	}

	return nil
}

func (uc *customOptionUsecase) DeleteCustomOption(ctx context.Context, id string) error {
	if err := uc.repo.DeleteCustomOption(ctx, id); err != nil {
		return fmt.Errorf("failed to delete custom option - %w", err)
	}

	return nil
}
