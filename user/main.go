package main

import (
	"context"
	"net/http"
	"user/controllers"
	"user/dal"
	_ "user/middlewares"
)

func main() {
	conString := "user=admin password=admin dbname=chat host=db port=5432 sslmode=disable"
	db := db.Connect(conString)

	dbMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	http.Handle("/register", dbMiddleware(http.HandlerFunc(controllers.RegisterController)))
	http.Handle("/login", dbMiddleware(http.HandlerFunc(controllers.SigninController)))

	http.ListenAndServe(":8080", nil)
}
