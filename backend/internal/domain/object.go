package domain

import (
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
	ObjectCustomOptions []ObjectCustomOption
}

type ObjectFilter struct {
	Limit        int
	Offset       int
	OrderBy      string
	Name         string
	ComparisonId string
}

var providedObjectOrderings = []string{"created_at", "name", "rating"}

func NewObjectFilter(limit, offset int, orderBy, name, comparisonId string) (ObjectFilter, error) {
	if offset < 0 || limit < 0 {
		return ObjectFilter{}, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	if !slices.Contains(providedObjectOrderings, orderBy) {
		return ObjectFilter{}, fmt.Errorf("incorrect ordering value")
	}

	return ObjectFilter{
		Limit:   limit,
		Offset:  offset,
		Name:    name,
		OrderBy: orderBy,
	}, nil
}
