package httptransport

import (
	"net/http"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
)

func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func (a *API) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			writeUnauthorized(w)
			return
		}

		principal, err := a.tokens.Parse(token)
		if err != nil {
			writeUnauthorized(w)
			return
		}

		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), principal)))
	})
}

func (a *API) requireRole(role commondomain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				writeUnauthorized(w)
				return
			}

			if principal.Role != role {
				writeForbidden(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
