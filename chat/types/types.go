package types

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SendMsg struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	Destination int `json:"destination"`
}

type Message struct {
    Type string `json:"type"`
    Message string `json:"message"`
    Sender int `json:"sender"`
    IsRead bool `json:"isRead"`
    ChatId int `json:"chatid"`
}

type UnreadMessagesReq struct {
    Type string `json:"type"`
    ChatId int `json:"chatid"`
}


type AllMessagesReq struct {
    Type string `json:"type"`
    ChatId int `json:"chatid"`
}

type AuthClaim struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

type MessageDB struct {
    Id int
    UserId int
    Message string
    IsRead bool
    CreatedAt time.Time
}

type NewMessage struct {
    Type string `json:"type"`
    ChatId int `json:"chatId"`
}
