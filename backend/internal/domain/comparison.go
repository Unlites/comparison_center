package domain

import (
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

func NewComparisonFilter(limit, offset int, orderBy string) (ComparisonFilter, error) {
	if offset < 0 || limit < 0 {
		return ComparisonFilter{}, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	if !slices.Contains(providedComparisonOrderings, orderBy) {
		return ComparisonFilter{}, fmt.Errorf("incorrect ordering value")
	}

	return ComparisonFilter{
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy,
	}, nil
}
