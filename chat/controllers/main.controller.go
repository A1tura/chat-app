package controllers

import (
	"chat/types"
	"chat/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	godotenv.Load()

	con, mgr, kafkaConnections, _, errorThrower := utils.LoadContext(r.Context())

	var id int

	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		fmt.Print("ff")
		errorThrower.ThrowError(403, "Invalid JWT token", true)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &types.AuthClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		errorThrower.ThrowError(403, "Invalid JWT token", true)
		return
	}

	if claims, ok := token.Claims.(*types.AuthClaim); ok && token.Valid {
		id, err = strconv.Atoi(claims.Id)
		if err != nil {
			errorThrower.ThrowError(503, "Internal Server Error", true)
			return
		}
		errorThrower.UserId = id
		mgr.Store(id, con)
		kafkaConnections.SendConnection(id)
	} else {
		errorThrower.ThrowError(403, "Invalid JWT token", true)
		return
	}

	for {
		var rawMessage map[string]interface{}

		if err := con.ReadJSON(&rawMessage); err != nil {
			errorThrower.ThrowError(404, "xdd", true)
			kafkaConnections.SendDisconnection(id)
			return
		}

		messageType := rawMessage["type"].(string)

		switch messageType {
		case "sendMessage":
			var v types.SendMsg
			if err := utils.MapToStruct(rawMessage, &v); err != nil {
				panic(err)
			}
			sendMessage(v, id, r.Context())
		case "Ping":
			Ping(r.Context())
		case "getUnreadMessages":
			var unreadMessagesReq types.UnreadMessagesReq

			if err := utils.MapToStruct(rawMessage, &unreadMessagesReq); err != nil {
				panic(err)
			}

			getUnreadMessages(unreadMessagesReq, r.Context(), id)

		case "getAllMessages":
			var allMessagesReq types.AllMessagesReq

			if err := utils.MapToStruct(rawMessage, &allMessagesReq); err != nil {
				panic(err)
			}

			getAllMessages(allMessagesReq, r.Context(), id)
		}
	}
}

func sendMessage(m types.SendMsg, sender int, ctx context.Context) {
	var chatId int64
	var chatExist bool
	var usersExist bool

	_, mgr, _, db, errorThrower := utils.LoadContext(ctx)

	usersExist = db.UserExist(sender, m.Destination)

	if !usersExist {
		errorThrower.ThrowError(404, "User do not exist", false)
		return
	}

	chatExist, chatId = db.ChatExist(m.Destination, sender)
	//receiverIsOnline := db.IsOnline(m.Destination)
	receiverIsOnline := mgr.ConSaved(m.Destination)

	if !chatExist {
		chatId = db.CreateChat(m.Destination, sender)
	}

	if receiverIsOnline {
		fmt.Print("online")
		receiver := mgr.Get(m.Destination)

		if receiver == nil {
			errorThrower.ThrowError(404, "Dest WebSocket is unvaliable", false)
		}

		message := types.NewMessage{
			Type:    "newMessage",
			ChatId:  int(chatId),
		}

		db.SaveMessage(int(chatId), sender, m.Content)
		receiver.WriteJSON(message)
	} else {
		db.SaveMessage(int(chatId), sender, m.Content)
	}

}

func getUnreadMessages(req types.UnreadMessagesReq, ctx context.Context, userId int) {
	con, _, _, db, errorThrower := utils.LoadContext(ctx)

	if chatExist := db.ChatExistById(req.ChatId); !chatExist {
		errorThrower.ThrowError(404, "Chat do not exist", false)
		return
	}

	if ownChat := db.IsUserAccessChat(userId, req.ChatId); !ownChat {
		errorThrower.ThrowError(404, "Invalid operation", false)
		return
	}

	messages := db.GetUnreadMessages(req.ChatId)

	for _, message := range messages {
		messageStruct := types.Message{
			Type:    "message",
			Message: message.Message,
			Sender:  message.UserId,
            IsRead: message.IsRead,
			ChatId:  req.ChatId,
		}

		db.MarkAsRead(message.Id)

		con.WriteJSON(messageStruct)
	}
}

func getAllMessages(req types.AllMessagesReq, ctx context.Context, userId int) {
	con, _, _, db, errorThrower := utils.LoadContext(ctx)

	if chatExist := db.ChatExistById(req.ChatId); !chatExist {
		errorThrower.ThrowError(404, "Chat do not exist", false)
		return
	}

	if ownChat := db.IsUserAccessChat(userId, req.ChatId); !ownChat {
		errorThrower.ThrowError(404, "Invalid operation", false)
		return
	}

	messages := db.GetAllMessages(req.ChatId)

	for _, message := range messages {
		messageStruct := types.Message{
			Type:    "message",
			Message: message.Message,
			Sender:  message.UserId,
            IsRead: message.IsRead,
			ChatId:  req.ChatId,
		}

        if !message.IsRead && message.UserId == userId {
            db.MarkAsRead(message.Id)
        }

		db.MarkAsRead(message.Id)

		con.WriteJSON(messageStruct)
	}
}

func Ping(ctx context.Context) {
	con, _, _, _, _ := utils.LoadContext(ctx)

	con.WriteMessage(1, []byte(string("Pong")))
}
