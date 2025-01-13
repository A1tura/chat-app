package error

import (
	"chat/kafka"
	"chat/manager"

	"github.com/gorilla/websocket"
)

type ErrorThrower struct {
	Mgr    *manager.Manager
	UserId int
	WS     *websocket.Conn
	kafka  *kafka.KafkaManagerConnections
}

type Error struct {
	Types    string `json:"type"`
	Code     int    `json:"code"`
	Content  string `json:"error"`
	Critical bool   `json:"critical"`
}

func (e *ErrorThrower) ThrowError(code int, content string, critical bool) {
	error := Error{
		Types:    "error",
		Code:     code,
		Content:  content,
		Critical: critical,
	}

	e.WS.WriteJSON(error)

	if critical {
		e.WS.Close()
		conSaved := e.Mgr.ConSaved(e.UserId)

		if conSaved {
			e.kafka.SendDisconnection(e.UserId)
			e.Mgr.Delete(e.UserId)
		}

		return
	}
}
