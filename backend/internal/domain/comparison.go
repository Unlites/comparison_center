package domain

import (
	"context"
	"fmt"
	"slices"
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
	OrderBy string
}

var providedComparisonOrderings = []string{"created_at"}

func NewComparisonFilter(limit, offset int, orderBy string) (*ComparisonFilter, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	if !slices.Contains(providedComparisonOrderings, orderBy) {
		return nil, fmt.Errorf("incorrect ordering value")
	}

	return &ComparisonFilter{
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy,
	}, nil
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
