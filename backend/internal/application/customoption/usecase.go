package customoption

import (
	"context"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type CustomOptionUsecase struct {
	repo      CustomOptionRepository
	generator IdGenerator
}

type CustomOptionRepository interface {
	GetCustomOptions(ctx context.Context, filter domain.CustomOptionFilter) ([]domain.CustomOption, error)
	GetCustomOptionById(ctx context.Context, id string) (domain.CustomOption, error)
	UpdateCustomOption(ctx context.Context, customOption domain.CustomOption) error
	CreateCustomOption(ctx context.Context, customOption domain.CustomOption) error
	DeleteCustomOption(ctx context.Context, id string) error
}

type IdGenerator interface {
	GenerateId() string
}

func NewCustomOptionUsecase(
	repo CustomOptionRepository,
	generator IdGenerator,
) *CustomOptionUsecase {
	return &CustomOptionUsecase{repo: repo, generator: generator}
}

func (uc *CustomOptionUsecase) GetCustomOptions(
	ctx context.Context,
	filter domain.CustomOptionFilter,
) ([]domain.CustomOption, error) {
	customOptions, err := uc.repo.GetCustomOptions(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom options - %w", err)
	}

	return customOptions, nil
}

func (uc *CustomOptionUsecase) GetCustomOptionById(
	ctx context.Context,
	id string,
) (domain.CustomOption, error) {
	customOption, err := uc.repo.GetCustomOptionById(ctx, id)
	if err != nil {
		return domain.CustomOption{}, fmt.Errorf("failed to get custom option - %w", err)
	}

	return customOption, nil
}

func (uc *CustomOptionUsecase) UpdateCustomOption(
	ctx context.Context,
	id string,
	customOption domain.CustomOption,
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

func (uc *CustomOptionUsecase) CreateCustomOption(
	ctx context.Context,
	customOption domain.CustomOption,
) error {
	customOption.Id = uc.generator.GenerateId()
	if err := uc.repo.CreateCustomOption(ctx, customOption); err != nil {
		return fmt.Errorf("failed to create custom option - %w", err)
	}

	return nil
}

func (uc *CustomOptionUsecase) DeleteCustomOption(ctx context.Context, id string) error {
	if err := uc.repo.DeleteCustomOption(ctx, id); err != nil {
		return fmt.Errorf("failed to delete custom option - %w", err)
	}

	return nil
}
