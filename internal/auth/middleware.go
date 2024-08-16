package auth

import (
	"caa-test/internal/api/resp"
	"net/http"
)

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) *middleware {
	return &middleware{
		secretKey: secretKey,
	}
}

func (m *middleware) StaticToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenStr := r.Header.Get("Authorization")
		if tokenStr != m.secretKey {
			resp.WriteJSON(w, http.StatusUnauthorized, resp.HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
