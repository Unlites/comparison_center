package domain

import (
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

func NewCustomOptionFilter(limit, offset int, name string) (CustomOptionFilter, error) {
	if offset < 0 || limit < 0 {
		return CustomOptionFilter{}, fmt.Errorf("offset amd limit must not be less than zero")
	}

	if limit == 0 {
		limit = 10
	}

	return CustomOptionFilter{
		Limit:  limit,
		Offset: offset,
		Name:   name,
	}, nil
}
