package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Unlites/comparison_center/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ObjectRepositoryMongo struct {
	objectsColl *mongo.Collection
}

func NewObjectRepositoryMongo(client *mongo.Client) *ObjectRepositoryMongo {
	return &ObjectRepositoryMongo{
		objectsColl: client.Database("database").Collection("objects"),
	}
}

type objectMongo struct {
	Id           string    `bson:"_id"`
	Name         string    `bson:"name"`
	Rating       int       `bson:"rating"`
	CreatedAt    time.Time `bson:"created_at"`
	Advs         string    `bson:"advs"`
	Disadvs      string    `bson:"disadvs"`
	PhotoPath    string    `bson:"photo_path"`
	ComparisonId string    `bson:"comparison_id"`
}

func (repo *ObjectRepositoryMongo) GetObjects(
	ctx context.Context,
	filter domain.ObjectFilter,
) ([]domain.Object, error) {
	opts := options.Find().
		SetSort(bson.M{filter.OrderBy: 1}).
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))

	condition := bson.M{}

	if filter.Name != "" {
		condition["name"] = bson.M{
			"$regex":   filter.Name,
			"$options": "i",
		}
	}

	if filter.ComparisonId != "" {
		condition["comparison_id"] = filter.ComparisonId
	}

	cur, err := repo.objectsColl.Find(ctx, condition, opts)
	if err != nil {
		return nil, fmt.Errorf("fetch objects from mongo error: %w", err)
	}

	objects := make([]domain.Object, 0, filter.Limit)
	for cur.Next(ctx) {
		var obj objectMongo
		if err := cur.Decode(obj); err != nil {
			return nil, fmt.Errorf("decode mongo result error %w", err)
		}

		objects = append(objects, toDomainObject(obj))
	}

	return objects, nil
}

func (repo *ObjectRepositoryMongo) GetObjectById(
	ctx context.Context,
	id string,
) (domain.Object, error) {
	res := repo.objectsColl.FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return domain.Object{}, fmt.Errorf("object %w", domain.ErrNotFound)
		}

		return domain.Object{}, fmt.Errorf("get object from mongo error %w", res.Err())
	}

	var obj objectMongo
	if err := res.Decode(obj); err != nil {
		return domain.Object{}, fmt.Errorf("decode mongo result error %w", err)
	}

	return toDomainObject(obj), nil
}

func (repo *ObjectRepositoryMongo) CreateObject(
	ctx context.Context,
	object domain.Object,
) error {
	_, err := repo.objectsColl.InsertOne(ctx, toObjectMongo(object))
	if err != nil {
		return fmt.Errorf("insert to mongo error: %w", err)
	}

	return nil
}

func (repo *ObjectRepositoryMongo) UpdateObject(
	ctx context.Context,
	object domain.Object,
) error {
	res, err := repo.objectsColl.UpdateOne(
		ctx,
		bson.M{"_id": object.Id},
		bson.M{"$set": toObjectMongo(object)},
	)
	if err != nil {
		return fmt.Errorf("update at mongo error: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("object %w", domain.ErrNotFound)
	}

	return nil
}

func (repo *ObjectRepositoryMongo) DeleteObject(
	ctx context.Context,
	id string,
) error {
	res, err := repo.objectsColl.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete from mongo error: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("object %w", domain.ErrNotFound)
	}

	return nil
}

func toDomainObject(objMongo objectMongo) domain.Object {
	return domain.Object{
		Id:           objMongo.Id,
		Name:         objMongo.Name,
		Rating:       objMongo.Rating,
		CreatedAt:    objMongo.CreatedAt,
		Advs:         objMongo.Advs,
		Disadvs:      objMongo.Disadvs,
		PhotoPath:    objMongo.PhotoPath,
		ComparisonId: objMongo.ComparisonId,
	}
}

func toObjectMongo(obj domain.Object) objectMongo {
	return objectMongo{
		Id:           obj.Id,
		Name:         obj.Name,
		Rating:       obj.Rating,
		CreatedAt:    obj.CreatedAt,
		Advs:         obj.Advs,
		Disadvs:      obj.Disadvs,
		PhotoPath:    obj.PhotoPath,
		ComparisonId: obj.ComparisonId,
	}
}
