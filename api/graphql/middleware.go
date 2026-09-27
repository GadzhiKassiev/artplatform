package graphql

import (
	"net/http"
	"strings"

	authpb "artplatform/backend/proto/auth"
)

type AuthMiddleware struct {
	AuthClient authpb.AuthServiceClient
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			next.ServeHTTP(w, r)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := m.AuthClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{Token: token})
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := withUserID(r.Context(), resp.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
