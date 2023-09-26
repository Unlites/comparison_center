package domain

import "context"

type CustomOption struct {
	Id    string
	Name  string
	Value string
}

type CustomOptionFilter struct {
	Limit  int
	Offset int
	Name   string
}

type CustomOptionUsecase interface {
	GetCustomOptions(ctx context.Context, filter *CustomOptionFilter) ([]*CustomOption, error)
	GetCustomOptionById(ctx context.Context, id string) (*CustomOption, error)
	UpdateCustomOption(ctx context.Context, id string, customOption *CustomOption) error
	CreateCustomOption(ctx context.Context, customOption *CustomOption) (*CustomOption, error)
}
