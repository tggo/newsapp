package model

import "time"

type Post struct {
	ID        int64
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Filter struct {
	Start int
	Limit int
}

func (f *Filter) GetStart() int {
	if f.Start <= 0 {
		return 0
	}

	return f.Start
}
func (f *Filter) GetLimit() int {
	if f.Limit > 50 {
		return 50
	}
	if f.Limit <= 0 {
		return 10
	}
	return f.Limit
}
