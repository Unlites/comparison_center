package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
)

type objectUsecase struct {
	objRepo        ObjectRepository
	custOptObjRepo ObjectCustomOptionRepository
	generator      IdGenerator
}

type ObjectUsecase interface {
	GetObjects(ctx context.Context, filter domain.ObjectFilter) ([]domain.Object, error)
	GetObjectById(ctx context.Context, id string) (domain.Object, error)
	UpdateObject(ctx context.Context, id string, object domain.Object) error
	CreateObject(ctx context.Context, object domain.Object) (string, error)
	DeleteObject(ctx context.Context, id string) error
	SetObjectPhotoPath(ctx context.Context, id, path string) error
}

type ObjectRepository interface {
	GetObjects(ctx context.Context, filter domain.ObjectFilter) ([]domain.Object, error)
	GetObjectById(ctx context.Context, id string) (domain.Object, error)
	UpdateObject(ctx context.Context, object domain.Object) error
	CreateObject(ctx context.Context, object domain.Object) error
	DeleteObject(ctx context.Context, id string) error
}

type ObjectCustomOptionRepository interface {
	GetObjectCustomOptionsByObjectId(ctx context.Context, objectId string) ([]domain.ObjectCustomOption, error)
	AddObjectCustomOption(ctx context.Context, objectCustomOption domain.ObjectCustomOption) error
	UpdateObjectCustomOption(ctx context.Context, objectCustomOption domain.ObjectCustomOption) error
}

type IdGenerator interface {
	GenerateId() string
}

func NewObjectUsecase(
	objRepo ObjectRepository,
	custOptObjRepo ObjectCustomOptionRepository,
	generator IdGenerator,
) *objectUsecase {
	return &objectUsecase{
		objRepo:        objRepo,
		custOptObjRepo: custOptObjRepo,
		generator:      generator,
	}
}

func (uc *objectUsecase) GetObjects(
	ctx context.Context,
	filter domain.ObjectFilter,
) ([]domain.Object, error) {
	objects, err := uc.objRepo.GetObjects(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects - %w", err)
	}

	for i, obj := range objects {
		options, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, obj.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get custom options - %w", err)
		}

		objects[i].ObjectCustomOptions = options
	}

	return objects, nil
}

func (uc *objectUsecase) GetObjectById(
	ctx context.Context,
	id string,
) (domain.Object, error) {
	object, err := uc.objRepo.GetObjectById(ctx, id)
	if err != nil {
		return domain.Object{}, fmt.Errorf("failed to get object - %w", err)
	}

	options, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, object.Id)
	if err != nil {
		return domain.Object{}, fmt.Errorf("failed to get custom options - %w", err)
	}

	object.ObjectCustomOptions = options

	return object, nil
}

func (uc *objectUsecase) UpdateObject(
	ctx context.Context,
	id string,
	inputObject domain.Object,
) error {
	existingObject, err := uc.objRepo.GetObjectById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing object - %w", err)
	}

	inputObject.Id = existingObject.Id
	inputObject.CreatedAt = existingObject.CreatedAt
	inputObject.ComparisonId = existingObject.ComparisonId
	inputObject.PhotoPath = existingObject.PhotoPath

	if err := uc.objRepo.UpdateObject(ctx, inputObject); err != nil {
		return fmt.Errorf("failed to update object - %w", err)
	}

	existingObjectOptions, err := uc.custOptObjRepo.GetObjectCustomOptionsByObjectId(ctx, existingObject.Id)
	if err != nil {
		return fmt.Errorf("failed to get existing custom options - %w", err)
	}

	for i := range inputObject.ObjectCustomOptions {
		inputObject.ObjectCustomOptions[i].ObjectId = existingObject.Id

		if slices.ContainsFunc(existingObjectOptions, func(option domain.ObjectCustomOption) bool {
			return option.CustomOptionId == inputObject.ObjectCustomOptions[i].CustomOptionId
		}) {
			err := uc.custOptObjRepo.UpdateObjectCustomOption(ctx, inputObject.ObjectCustomOptions[i])
			if err != nil {
				return fmt.Errorf("failed to update custom option - %w", err)
			}
		} else {
			err := uc.custOptObjRepo.AddObjectCustomOption(ctx, inputObject.ObjectCustomOptions[i])
			if err != nil {
				return fmt.Errorf("failed to add custom option - %w", err)
			}
		}
	}

	return nil
}

func (uc *objectUsecase) CreateObject(
	ctx context.Context,
	object domain.Object,
) (string, error) {
	object.Id = uc.generator.GenerateId()
	object.CreatedAt = time.Now()
	if err := uc.objRepo.CreateObject(ctx, object); err != nil {
		return "", fmt.Errorf("failed to create object - %w", err)
	}

	return object.Id, nil
}

func (uc *objectUsecase) DeleteObject(ctx context.Context, id string) error {
	if err := uc.objRepo.DeleteObject(ctx, id); err != nil {
		return fmt.Errorf("failed to delete object - %w", err)
	}

	return nil
}

func (uc *objectUsecase) SetObjectPhotoPath(ctx context.Context, id, path string) error {
	object, err := uc.objRepo.GetObjectById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get object - %w", err)
	}

	object.PhotoPath = path

	if err := uc.objRepo.UpdateObject(ctx, object); err != nil {
		return fmt.Errorf("failed to update object - %w", err)
	}

	return nil
}
