package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"user/controllers"
	"user/dal"
	"user/kafka"
	_ "user/middlewares"

	"github.com/joho/godotenv"
)

func main() {
    godotenv.Load()
	var db_host string
	var kafka_host string

    db_host = os.Getenv("DB_HOST")
    kafka_host = os.Getenv("KAFKA_HOST")
	fmt.Printf("Host: %s", db_host)
	fmt.Printf("Kafka host: %s", kafka_host)
	conString := "user=admin password=admin dbname=chat host=" + db_host + " port=5432 sslmode=disable"
	db := db.Connect(conString)
	kafka := kafka.NewKafkaConnection([]string{kafka_host})

	go kafka.StartKafkaService(db)

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
