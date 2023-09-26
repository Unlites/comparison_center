package domain

import (
	"context"
	"time"
)

type Object struct {
	Id            string
	Name          string
	Rating        int
	CreatedAt     time.Time
	Advs          string
	Disadvs       string
	PhotoPath     string
	ComparisonId  string
	CustomOptions []*CustomOption
}

type ObjectFilter struct {
	Limit        int
	Offset       int
	OrderBy      string
	Name         string
	CustomOption CustomOption
}

type ObjectUsecase interface {
	GetObjects(ctx context.Context, filter *ObjectFilter) ([]*Object, error)
	GetObjectById(ctx context.Context, id string) (*Object, error)
	UpdateObject(ctx context.Context, id string, object *Object) error
	CreateObject(ctx context.Context, object *Object) (*Object, error)
}
