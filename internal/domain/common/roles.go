package common

import (
	"errors"
	"fmt"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

var ErrInvalidRole = errors.New("invalid role")

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

func (r Role) String() string {
	return string(r)
}

func ParseRole(raw string) (Role, error) {
	role := Role(raw)
	if !role.IsValid() {
		return "", fmt.Errorf("%w: %q", ErrInvalidRole, raw)
	}

	return role, nil
}
