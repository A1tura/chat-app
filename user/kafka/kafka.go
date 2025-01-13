package kafka

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"user/dal"

	"github.com/segmentio/kafka-go"
)

type KafkaConnections struct {
	Kafka_con *kafka.Reader
}

func NewKafkaConnection(Brokers []string) *KafkaConnections {
    fmt.Printf("\nBrokers: %s\n", Brokers)
	con := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  Brokers,
		GroupID:  "Controller",
		Topic:    "connections",
		MaxBytes: 10e6,
	})

	return &KafkaConnections{
		Kafka_con: con,
	}
}

func (kafka *KafkaConnections) StartKafkaService(db *db.DB) {
	for {
		msg, err := kafka.Kafka_con.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		id, err := strconv.Atoi(strings.Split(string(msg.Value), "-")[1])

		if err != nil {
			panic(err)
		}

		if strings.Split(string(msg.Value), "-")[0] == "connected" {
            fmt.Print(id)
			db.ChangeStatus(id, true)
		} else {
            fmt.Printf("dis: ", id)
            db.ChangeStatus(id, false)
        }
	}

}
