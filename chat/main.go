package main

import (
	"chat/controllers"
	"chat/dal"
	"chat/error"
	"chat/kafka"
	"chat/manager"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var mgr = manager.NewManager()

func main() {
	godotenv.Load()
	db_host := os.Getenv("DB_HOST")
    kafka_host := os.Getenv("KAFKA_HOST")
	db := db.NewDB("user=admin password=admin dbname=chat host=" + db_host + " port=5432 sslmode=disable")
	var kafkaConnections = kafka.NewKafkaManagerConnections(kafka_host)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	midl := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			con, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				con.WriteMessage(1, []byte("Error: Failed to upgrade to WebSocket"))
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, "ws", con)
			ctx = context.WithValue(ctx, "mgr", mgr)
			ctx = context.WithValue(ctx, "kafka", kafkaConnections)
			ctx = context.WithValue(ctx, "db", db)
			ctx = context.WithValue(ctx, "errors", error.ErrorThrower{Mgr: mgr, WS: con, UserId: -1})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	http.Handle("/", midl(http.HandlerFunc(controllers.Handler)))

	http.ListenAndServe(":1337", nil)

	<-stop
	log.Println("Shutting the server")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Exec("UPDATE users SET isonline=false WHERE isonline=true")
	if err != nil {
		log.Printf("Error during shutting the server: %s", err)
	}
}
