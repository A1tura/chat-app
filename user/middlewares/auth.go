package middlewares

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"user/types"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Forbidden", 403)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &types.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			http.Error(w, "Forbidden", 403)
			return
		}

		if claims, ok := token.Claims.(*types.AuthClaims); ok && token.Valid {
			id, err := strconv.Atoi(claims.Id)
			if err != nil {
				http.Error(w, "External Error", 502)
				return
			}
			ctx := context.WithValue(r.Context(), "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Forbidden", 403)
			return
		}
	})
}
