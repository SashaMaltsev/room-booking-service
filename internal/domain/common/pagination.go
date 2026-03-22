package common

import (
	"errors"
	"fmt"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

var (
	ErrInvalidPage     = errors.New("page must be greater than zero")
	ErrInvalidPageSize = errors.New("page size is out of allowed range")
)

type PageRequest struct {
	Page     int
	PageSize int
}

func (r PageRequest) Normalize() PageRequest {
	if r.Page == 0 {
		r.Page = DefaultPage
	}

	if r.PageSize == 0 {
		r.PageSize = DefaultPageSize
	}

	return r
}

func (r PageRequest) Validate() error {
	normalized := r.Normalize()

	if normalized.Page < 1 {
		return fmt.Errorf("%w: %d", ErrInvalidPage, normalized.Page)
	}

	if normalized.PageSize < 1 || normalized.PageSize > MaxPageSize {
		return fmt.Errorf("%w: %d", ErrInvalidPageSize, normalized.PageSize)
	}

	return nil
}

func (r PageRequest) Offset() int {
	normalized := r.Normalize()
	return (normalized.Page - 1) * normalized.PageSize
}

type Page struct {
	Page     int
	PageSize int
	Total    int
}
