package domain

import (
	"context"
	"time"
)

type Comparison struct {
	Id              string
	Name            string
	CreatedAt       time.Time
	CustomOptionIds []string
}

type ComparisonFilter struct {
	Limit   int
	Offset  int
	OrderBy int
}

type ComparisonUsecase interface {
	GetComparisons(ctx context.Context, filter *ComparisonFilter) ([]*Comparison, error)
	GetComparisonById(ctx context.Context, id string) (*Comparison, error)
	UpdateComparison(ctx context.Context, id string, comparison *Comparison) error
	CreateComparison(ctx context.Context, comparison *Comparison) error
	DeleteComparison(ctx context.Context, id string) error
}

type ComparisonRepository interface {
	GetComparisons(ctx context.Context, filter *ComparisonFilter) ([]*Comparison, error)
	GetComparisonById(ctx context.Context, id string) (*Comparison, error)
	UpdateComparison(ctx context.Context, comparison *Comparison) error
	CreateComparison(ctx context.Context, comparison *Comparison) error
	DeleteComparison(ctx context.Context, id string) error
}
