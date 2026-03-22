package common

import (
	"errors"
	"testing"
)

func TestPageRequestNormalizeUsesDefaults(t *testing.T) {
	request := PageRequest{}

	normalized := request.Normalize()

	if normalized.Page != DefaultPage {
		t.Fatalf("expected default page %d, got %d", DefaultPage, normalized.Page)
	}

	if normalized.PageSize != DefaultPageSize {
		t.Fatalf("expected default page size %d, got %d", DefaultPageSize, normalized.PageSize)
	}
}

func TestPageRequestValidateRejectsTooLargePageSize(t *testing.T) {
	request := PageRequest{Page: 1, PageSize: MaxPageSize + 1}

	err := request.Validate()
	if !errors.Is(err, ErrInvalidPageSize) {
		t.Fatalf("expected ErrInvalidPageSize, got %v", err)
	}
}
