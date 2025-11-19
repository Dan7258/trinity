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

func SoftAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if !token.Valid {
			next.ServeHTTP(w, r)
			return
		}
		claims, err := jwt.ParseClaims(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(jwt.Claims)
		if !ok {
			http.Error(w, "failed to get user claimss", http.StatusUnauthorized)
			return
		}
		if claims.Role != "admin" {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		//log.Printf("id: %d \n role: %s \n login: %s \n", claims.ID, claims.Role, claims.Login)
		next.ServeHTTP(w, r)
	})
}
