package middleware

import (
	"context"
	"fmt"
	"net/http"
	"trinity/pkg/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, `{"message": "bad token"}`, http.StatusUnauthorized)
			return
		}
		claims, err := jwt.ParseClaims(token)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"message": "%v"}`, err), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
