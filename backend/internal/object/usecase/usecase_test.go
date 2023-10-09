package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	or "github.com/Unlites/comparison_center/backend/internal/object/repository"
	cr "github.com/Unlites/comparison_center/backend/internal/object_customoption/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetObjects(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		returnedObjects := []*domain.Object{
			{
				Id:           "231934sadas9123deqw",
				Name:         "BMW X5",
				Rating:       8,
				CreatedAt:    time.Now(),
				Advs:         "Good SUV",
				Disadvs:      "Hard to find some details",
				PhotoPath:    "/cars/231934sadas9123deqw.jpg",
				ComparisonId: "85434230werhuhi123912304",
			},
		}

		returnedOptions := []*domain.ObjectCustomOption{
			{
				ObjectId:       "231934sadas9123deqw",
				CustomOptionId: "432230ewrew3424rwe",
				Value:          "600",
			},
			{
				ObjectId:       "231934sadas9123deqw",
				CustomOptionId: "52342rwerew23123",
				Value:          "2021",
			},
		}

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()
		filter := &domain.ObjectFilter{
			Limit:   2,
			Offset:  0,
			OrderBy: "created",
		}

		objRepo.On("GetObjects", ctx, filter).Return(returnedObjects, nil)
		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedObjects[0].Id).Return(returnedOptions, nil)

		objects, err := uc.GetObjects(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, objects, returnedObjects)
		assert.Equal(t, returnedOptions, returnedObjects[0].ObjectCustomOptions)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()
		filter := &domain.ObjectFilter{
			Limit:  2,
			Offset: 0,
		}

		objRepo.On("GetObjects", ctx, filter).Return(nil, errors.New("some error"))
		objects, err := uc.GetObjects(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, objects)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
		objRepo.AssertExpectations(t)
	})
}

func TestGetObjectById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		returnedObject := &domain.Object{

			Id:           "231934sadas9123deqw",
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		returnedOptions := []*domain.ObjectCustomOption{
			{
				ObjectId:       "231934sadas9123deqw",
				CustomOptionId: "432230ewrew3424rwe",
				Value:          "600",
			},
			{
				ObjectId:       "231934sadas9123deqw",
				CustomOptionId: "52342rwerew23123",
				Value:          "2021",
			},
		}

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()
		id := "231934sadas9123deqw"

		objRepo.On("GetObjectById", ctx, id).Return(returnedObject, nil)
		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedObject.Id).Return(returnedOptions, nil)

		object, err := uc.GetObjectById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, object, returnedObject)
		assert.Equal(t, returnedOptions, returnedObject.ObjectCustomOptions)

		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()
		id := "213213ewrwe9423432"

		objRepo.On("GetObjectById", ctx, id).Return(nil, domain.ErrNotFound)

		object, err := uc.GetObjectById(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, object)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
	})
}

func TestCreateObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		inputObject := &domain.Object{
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()

		objRepo.On("CreateObject", ctx, inputObject).Return(nil)

		err := uc.CreateObject(ctx, inputObject)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		inputObject := &domain.Object{
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()

		objRepo.On("CreateObject", ctx, inputObject).Return(errors.New("some error"))

		err := uc.CreateObject(ctx, inputObject)

		assert.Error(t, err)
		objRepo.AssertExpectations(t)
	})
}

func TestUpdateObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		returnedOnGetObject := &domain.Object{
			Id:           "231934sadas9123deqw",
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []*domain.ObjectCustomOption{
				{
					ObjectId:       "231934sadas9123deqw",
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "600",
				},
			},
		}

		inputObject := &domain.Object{
			Name:         "BMW X5",
			Rating:       9,
			CreatedAt:    time.Now(),
			Advs:         "Very good SUV",
			Disadvs:      "Easy to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []*domain.ObjectCustomOption{
				{
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "800",
				},
			},
		}

		id := "231934sadas9123deqw"

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(returnedOnGetObject, nil)
		objRepo.On("UpdateObject", ctx, inputObject).Return(nil)

		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedOnGetObject.Id).
			Return(returnedOnGetObject.ObjectCustomOptions, nil)

		custOptObjRepo.On("UpdateObjectCustomOption", ctx, inputObject.ObjectCustomOptions[0]).Return(nil)

		err := uc.UpdateObject(ctx, id, inputObject)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		inputObject := &domain.Object{
			Name:         "BMW X5",
			Rating:       9,
			CreatedAt:    time.Now(),
			Advs:         "Very good SUV",
			Disadvs:      "Easy to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []*domain.ObjectCustomOption{
				{
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "800",
				},
			},
		}

		id := "231934sadas9123deqw"

		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(nil, errors.New("some error"))

		err := uc.UpdateObject(ctx, id, inputObject)

		assert.Error(t, err)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
		objRepo.AssertNotCalled(t, "UpdateObject")
		objRepo.AssertNotCalled(t, "UpdateObjectCustomOption")
		custOptObjRepo.AssertExpectations(t)
	})
}

func TestDeleteObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)

		ctx := context.Background()
		id := "34543dfsdfj32432jewr"

		objRepo.On("DeleteObject", ctx, id).Return(nil)

		err := uc.DeleteObject(ctx, id)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := or.NewObjectRepositoryMock()
		custOptObjRepo := cr.NewObjectCustomOptionRepositoryMock()
		uc := NewObjectUsecase(objRepo, custOptObjRepo)
		ctx := context.Background()
		id := "92133easd123srewr132"

		objRepo.On("DeleteObject", ctx, id).Return(errors.New("some error"))

		err := uc.DeleteObject(ctx, id)

		assert.Error(t, err)
		objRepo.AssertExpectations(t)
	})
}
