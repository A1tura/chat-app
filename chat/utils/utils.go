package utils

import (
	db "chat/dal"
	customError "chat/error"
	"chat/kafka"
	"chat/manager"
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
)

func LoadContext(c context.Context) (*websocket.Conn, *manager.Manager, *kafka.KafkaManagerConnections, *db.DB, customError.ErrorThrower) {
	ws := c.Value("ws").(*websocket.Conn)
	mgr := c.Value("mgr").(*manager.Manager)
	kafka := c.Value("kafka").(*kafka.KafkaManagerConnections)
	db := c.Value("db").(*db.DB)
	errors := c.Value("errors").(customError.ErrorThrower)

	return ws, mgr, kafka, db, errors
}

func MapToStruct(data map[string]interface{}, out interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, out)
}



func UNUSED(x ...interface{}) {}
