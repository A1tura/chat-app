package kafka

import (
	"context"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type KafkaManagerConnections struct {
	kafka_con *kafka.Conn
}

func NewKafkaManagerConnections(host string) *KafkaManagerConnections {
	con, err := kafka.DialLeader(context.Background(), "tcp", host, "connections", 0)
	if err != nil {
		panic(err)
	}
	return &KafkaManagerConnections{
		kafka_con: con,
	}
}

func (kafka *KafkaManagerConnections) SendConnection(id int) {
	kafka.kafka_con.Write([]byte("connected-" + strconv.Itoa(id)))
	fmt.Print("ff")
}

func (kafka *KafkaManagerConnections) SendDisconnection(id int) {
	kafka.kafka_con.Write([]byte("disconnected-" + strconv.Itoa(id)))
}
