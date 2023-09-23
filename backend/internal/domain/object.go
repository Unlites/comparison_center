package domain

import "time"

type Object struct {
	Id            string
	Name          string
	Rating        int
	CreatedAt     time.Time
	Advs          string
	Disadvs       []string
	PhotoPath     string
	ComparisonId  string
	CustomOptions []*CustomOption
}
