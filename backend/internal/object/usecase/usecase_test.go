package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"github.com/Unlites/comparison_center/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetObjects(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)
		returnedObjects := []domain.Object{
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

		returnedOptions := []domain.ObjectCustomOption{
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

		ctx := context.Background()
		filter := domain.ObjectFilter{
			Limit:   2,
			Offset:  0,
			OrderBy: "created",
		}

		objRepo.On("GetObjects", ctx, filter).Return(returnedObjects, nil)
		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedObjects[0].Id).Return(returnedOptions, nil)

		objects, err := uc.GetObjects(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, returnedObjects, objects)
		assert.Equal(t, returnedOptions, objects[0].ObjectCustomOptions)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		ctx := context.Background()
		filter := domain.ObjectFilter{
			Limit:  2,
			Offset: 0,
		}

		objRepo.On("GetObjects", ctx, filter).Return(nil, assert.AnError)
		objects, err := uc.GetObjects(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, objects)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
		objRepo.AssertExpectations(t)
	})
}

func TestGetObjectById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		returnedObject := domain.Object{
			Id:           "231934sadas9123deqw",
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		returnedOptions := []domain.ObjectCustomOption{
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

		returnedObjectWithOptions := returnedObject
		returnedObjectWithOptions.ObjectCustomOptions = returnedOptions

		ctx := context.Background()
		id := "231934sadas9123deqw"

		objRepo.On("GetObjectById", ctx, id).Return(returnedObject, nil)
		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedObject.Id).Return(returnedOptions, nil)

		object, err := uc.GetObjectById(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, returnedObjectWithOptions, object)

		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		ctx := context.Background()
		id := "213213ewrwe9423432"

		objRepo.On("GetObjectById", ctx, id).Return(nil, domain.ErrNotFound)

		object, err := uc.GetObjectById(ctx, id)

		assert.Error(t, err)
		assert.Empty(t, object)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
	})
}

func TestCreateObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		inputObject := domain.Object{
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		ctx := context.Background()

		objRepo.On("CreateObject", ctx, mock.MatchedBy(func(object domain.Object) bool {
			return object.Name == inputObject.Name &&
				object.Rating == inputObject.Rating &&
				object.Advs == inputObject.Advs &&
				object.Disadvs == inputObject.Disadvs &&
				object.PhotoPath == inputObject.PhotoPath &&
				object.ComparisonId == inputObject.ComparisonId
		})).Return(nil)
		generator.On("GenerateId").Return("231934sadas9123deqw")

		id, err := uc.CreateObject(ctx, inputObject)

		assert.Equal(t, "231934sadas9123deqw", id)
		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		inputObject := domain.Object{
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
		}

		ctx := context.Background()

		objRepo.On("CreateObject", ctx, mock.MatchedBy(func(object domain.Object) bool {
			return object.Name == inputObject.Name &&
				object.Rating == inputObject.Rating &&
				object.Advs == inputObject.Advs &&
				object.Disadvs == inputObject.Disadvs &&
				object.PhotoPath == inputObject.PhotoPath &&
				object.ComparisonId == inputObject.ComparisonId
		})).Return(assert.AnError)
		generator.On("GenerateId").Return("231934sadas9123deqw")

		id, err := uc.CreateObject(ctx, inputObject)

		assert.Empty(t, id)
		assert.Error(t, err)
		objRepo.AssertExpectations(t)
	})
}

func TestUpdateObject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		returnedOnGetObject := domain.Object{
			Id:           "231934sadas9123deqw",
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []domain.ObjectCustomOption{
				{
					ObjectId:       "231934sadas9123deqw",
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "600",
				},
			},
		}

		inputObject := domain.Object{
			Name:         "BMW X5",
			Rating:       9,
			CreatedAt:    time.Now(),
			Advs:         "Very good SUV",
			Disadvs:      "Easy to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []domain.ObjectCustomOption{
				{
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "800",
				},
			},
		}

		changedObject := inputObject
		changedObject.Id = returnedOnGetObject.Id
		changedObject.ObjectCustomOptions[0].ObjectId = returnedOnGetObject.Id

		id := "231934sadas9123deqw"

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(returnedOnGetObject, nil)
		objRepo.On("UpdateObject", ctx, changedObject).Return(nil)

		custOptObjRepo.On("GetObjectCustomOptionsByObjectId", ctx, returnedOnGetObject.Id).
			Return(returnedOnGetObject.ObjectCustomOptions, nil)

		custOptObjRepo.On("UpdateObjectCustomOption", ctx, inputObject.ObjectCustomOptions[0]).Return(nil)

		err := uc.UpdateObject(ctx, id, inputObject)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
		custOptObjRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		inputObject := domain.Object{
			Name:         "BMW X5",
			Rating:       9,
			CreatedAt:    time.Now(),
			Advs:         "Very good SUV",
			Disadvs:      "Easy to find some details",
			PhotoPath:    "/cars/231934sadas9123deqw.jpg",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []domain.ObjectCustomOption{
				{
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "800",
				},
			},
		}

		id := "231934sadas9123deqw"

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(nil, assert.AnError)

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
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		ctx := context.Background()
		id := "34543dfsdfj32432jewr"

		objRepo.On("DeleteObject", ctx, id).Return(nil)

		err := uc.DeleteObject(ctx, id)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		ctx := context.Background()
		id := "92133easd123srewr132"

		objRepo.On("DeleteObject", ctx, id).Return(assert.AnError)

		err := uc.DeleteObject(ctx, id)

		assert.Error(t, err)
		objRepo.AssertExpectations(t)
	})
}

func TestSetObjectPhotoPath(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		object := domain.Object{
			Id:           "231934sadas9123deqw",
			Name:         "BMW X5",
			Rating:       8,
			CreatedAt:    time.Now(),
			Advs:         "Good SUV",
			Disadvs:      "Hard to find some details",
			ComparisonId: "85434230werhuhi123912304",
			ObjectCustomOptions: []domain.ObjectCustomOption{
				{
					ObjectId:       "231934sadas9123deqw",
					CustomOptionId: "432230ewrew3424rwe",
					Value:          "600",
				},
			},
		}

		id := "92133easd123srewr132"
		path := "/photos/4324123sfnjsadn1239213.jpg"

		changedObject := object
		changedObject.PhotoPath = path

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(object, nil)
		objRepo.On("UpdateObject", ctx, changedObject).Return(nil)

		err := uc.SetObjectPhotoPath(ctx, id, path)

		assert.NoError(t, err)
		objRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		objRepo := mocks.NewObjectRepositoryMock()
		custOptObjRepo := mocks.NewObjectCustomOptionRepositoryMock()
		generator := mocks.NewMockGenerator()
		uc := NewObjectUsecase(objRepo, custOptObjRepo, generator)

		id := "231934sadas9123deqw"
		path := "/photos/4324123sfnjsadn1239213.jpg"

		ctx := context.Background()

		objRepo.On("GetObjectById", ctx, id).Return(nil, assert.AnError)

		err := uc.SetObjectPhotoPath(ctx, id, path)

		assert.Error(t, err)
		custOptObjRepo.AssertNotCalled(t, "GetObjectCustomOptionsByObjectId")
		objRepo.AssertNotCalled(t, "UpdateObject")
	})
}
