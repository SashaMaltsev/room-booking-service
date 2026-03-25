package httptransport

import (
	"context"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

type Principal struct {
	UserID string
	Role   commondomain.Role
}

type contextKey string

const principalContextKey contextKey = "principal"

func WithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey).(Principal)
	return principal, ok
}
