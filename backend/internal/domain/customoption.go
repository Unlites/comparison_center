package domain

import (
	"context"
	"fmt"
)

type CustomOption struct {
	Id   string
	Name string
}

type CustomOptionFilter struct {
	Limit  int
	Offset int
	Name   string
}

func NewCustomOptionFilter(limit, offset int, name string) (*CustomOptionFilter, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	return &CustomOptionFilter{
		Limit:  limit,
		Offset: offset,
		Name:   name,
	}, nil
}

type CustomOptionUsecase interface {
	GetCustomOptions(ctx context.Context, filter *CustomOptionFilter) ([]*CustomOption, error)
	GetCustomOptionById(ctx context.Context, id string) (*CustomOption, error)
	UpdateCustomOption(ctx context.Context, id string, customOption *CustomOption) error
	CreateCustomOption(ctx context.Context, customOption *CustomOption) error
	DeleteCustomOption(ctx context.Context, id string) error
}

type CustomOptionRepository interface {
	GetCustomOptions(ctx context.Context, filter *CustomOptionFilter) ([]*CustomOption, error)
	GetCustomOptionById(ctx context.Context, id string) (*CustomOption, error)
	UpdateCustomOption(ctx context.Context, customOption *CustomOption) error
	CreateCustomOption(ctx context.Context, customOption *CustomOption) error
	DeleteCustomOption(ctx context.Context, id string) error
}
