package auth

import (
	"errors"
	"testing"
	"time"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

func TestIssueAndParse(t *testing.T) {
	manager := NewManager("secret", time.Hour)
	manager.now = func() time.Time {
		return time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
	}

	token, err := manager.Issue("user-1", commondomain.RoleUser)
	if err != nil {
		t.Fatalf("expected token issue without error, got %v", err)
	}

	principal, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("expected token parse without error, got %v", err)
	}

	if principal.UserID != "user-1" {
		t.Fatalf("expected user id user-1, got %q", principal.UserID)
	}

	if principal.Role != commondomain.RoleUser {
		t.Fatalf("expected role user, got %q", principal.Role)
	}
}

func TestParseExpiredToken(t *testing.T) {
	manager := NewManager("secret", time.Hour)
	manager.now = func() time.Time {
		return time.Date(2026, time.March, 25, 10, 0, 0, 0, time.UTC)
	}

	token, err := manager.Issue("user-1", commondomain.RoleUser)
	if err != nil {
		t.Fatalf("expected token issue without error, got %v", err)
	}

	manager.now = func() time.Time {
		return time.Date(2026, time.March, 25, 12, 0, 1, 0, time.UTC)
	}

	_, err = manager.Parse(token)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}
