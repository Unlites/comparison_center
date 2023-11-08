package domain

import (
	"context"
	"fmt"
	"slices"
	"time"
)

type Object struct {
	Id                  string
	Name                string
	Rating              int
	CreatedAt           time.Time
	Advs                string
	Disadvs             string
	PhotoPath           string
	ComparisonId        string
	ObjectCustomOptions []*ObjectCustomOption
}

type ObjectFilter struct {
	Limit   int
	Offset  int
	OrderBy string
	Name    string
}

func NewObjectFilter(limit, offset int, orderBy, name string) (*ObjectFilter, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	providedOrderings := []string{"created_at", "name", "rating"}

	if !slices.Contains(providedOrderings, orderBy) {
		return nil, fmt.Errorf("incorrect ordering value")
	}

	return &ObjectFilter{
		Limit:   limit,
		Offset:  offset,
		Name:    name,
		OrderBy: orderBy,
	}, nil
}

type ObjectUsecase interface {
	GetObjects(ctx context.Context, filter *ObjectFilter) ([]*Object, error)
	GetObjectById(ctx context.Context, id string) (*Object, error)
	UpdateObject(ctx context.Context, id string, object *Object) error
	CreateObject(ctx context.Context, object *Object) error
	DeleteObject(ctx context.Context, id string) error
	SetObjectPhotoPath(ctx context.Context, id, path string) error
}

type ObjectRepository interface {
	GetObjects(ctx context.Context, filter *ObjectFilter) ([]*Object, error)
	GetObjectById(ctx context.Context, id string) (*Object, error)
	UpdateObject(ctx context.Context, object *Object) error
	CreateObject(ctx context.Context, object *Object) error
	DeleteObject(ctx context.Context, id string) error
}
