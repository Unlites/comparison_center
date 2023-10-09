package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type objectUsecase struct {
	objRepo        domain.ObjectRepository
	custOptObjRepo domain.ObjectCustomOptionRepository
}

func NewObjectUsecase(
	objRepo domain.ObjectRepository,
	custOptObjRepo domain.ObjectCustomOptionRepository,
) *objectUsecase {
	return &objectUsecase{objRepo: objRepo, custOptObjRepo: custOptObjRepo}
}

func (uc *objectUsecase) GetObjects(
	ctx context.Context,
	filter *domain.ObjectFilter,
) ([]*domain.Object, error) {
	objects, err := uc.objRepo.GetObjects(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects - %w", err)
	}

	for _, obj := range objects {
		options, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, obj.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get custom options - %w", err)
		}

		obj.ObjectCustomOptions = options
	}

	return objects, nil
}

func (uc *objectUsecase) GetObjectById(
	ctx context.Context,
	id string,
) (*domain.Object, error) {
	object, err := uc.objRepo.GetObjectById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get object - %w", err)
	}

	options, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, object.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom options - %w", err)
	}

	object.ObjectCustomOptions = options

	return object, nil
}

func (uc *objectUsecase) UpdateObject(
	ctx context.Context,
	id string,
	inputObject *domain.Object,
) error {
	existingObject, err := uc.objRepo.GetObjectById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing object - %w", err)
	}

	if err := uc.objRepo.UpdateObject(ctx, inputObject); err != nil {
		return fmt.Errorf("failed to update object - %w", err)
	}

	for _, option := range inputObject.ObjectCustomOptions {
		existingObjectOption, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, existingObject.Id)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("failed to get existing custom option - %w", err)
		}

		if existingObjectOption == nil {
			err := uc.custOptObjRepo.AddObjectCustomOption(ctx, option)
			if err != nil {
				return fmt.Errorf("failed to add custom option - %w", err)
			}
		} else {
			err := uc.custOptObjRepo.UpdateObjectCustomOption(ctx, option)
			if err != nil {
				return fmt.Errorf("failed to update custom option - %w", err)
			}
		}
	}

	return nil
}

func (uc *objectUsecase) CreateObject(
	ctx context.Context,
	object *domain.Object,
) error {
	if err := uc.objRepo.CreateObject(ctx, object); err != nil {
		return fmt.Errorf("failed to create object - %w", err)
	}

	return nil
}

func (uc *objectUsecase) DeleteObject(ctx context.Context, id string) error {
	if err := uc.objRepo.DeleteObject(ctx, id); err != nil {
		return fmt.Errorf("failed to delete object - %w", err)
	}

	return nil
}
