package db

import (
	"chat/array"
	"chat/types"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(host string) *DB {
	db, err := sql.Open("postgres", host)

	if err != nil {
		panic(err)
	}

	if db.Ping(); err != nil {
		panic(err)
	}

	return &DB{
		DB: db,
	}
}

func (db *DB) ChatExistById(chatId int) bool {
	var exist int
	row := db.DB.QueryRow(`SELECT id FROM "chat" WHERE id = $1`, chatId)

	row.Scan(&exist)
	fmt.Println(exist)

	if chatId == exist {
		return true
	}

	return false

}

func (db *DB) ChatExist(userA, userB int) (bool, int64) {
	query := `SELECT id FROM "chat" WHERE users = $1::bigint[];`
	row := db.DB.QueryRow(query, pq.Array([]int64{int64(userA), int64(userB)}))

	var id int64
	err := row.Scan(&id) // Only call Scan if a row exists
	if err != nil {
		if err == sql.ErrNoRows {
			query := `SELECT id FROM "chat" WHERE users = $1::bigint[];`
			row := db.DB.QueryRow(query, pq.Array([]int64{int64(userB), int64(userA)}))

			err := row.Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					return false, id
				}
			}
		}
		panic(err) // Handle unexpected errors
	}

	return true, id
}

func (db *DB) IsUserAccessChat(user int, chat int) bool {
	var users array.Array[int]
	var userResArray pq.Int64Array
	row := db.DB.QueryRow(`SELECT users FROM "chat" WHERE id = $1`, chat)

	row.Scan(&userResArray)

	for _, u := range userResArray {
		users = append(users, int(u))
	}

	exist, _ := users.ArrayInclude(user)

	return exist
}

func (db *DB) CreateChat(userA, userB int) int64 {
	var id int64
	row, err := db.DB.Exec(`INSERT INTO "chat" (users) VALUES ($1::bigint[]);`, pq.Array([]int64{int64(userA), int64(userB)}))
	if err != nil {
		panic(err)
	}

	idReq, err := row.LastInsertId()
	if err != nil {
		panic(err)
	}

	id = idReq

	return id
}

func (db *DB) GetUnreadMessages(chatId int) []types.MessageDB {
	var messages []types.MessageDB

	rows, err := db.DB.Query(`SELECT id, userId, message, isRead, createdAt FROM "message" WHERE chatId = $1 AND isread = false`, chatId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var messageStruct types.MessageDB
		var id int
		var userId int
		var message string
		var isRead bool
		var createdAt time.Time

		rows.Scan(&id, &userId, &message, &isRead, &createdAt)

		messageStruct = types.MessageDB{
			Id:        id,
			UserId:    userId,
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
		}

		messages = append(messages, messageStruct)
	}

	return messages
}

func (db *DB) GetAllMessages(chatId int) []types.MessageDB {
	var messages []types.MessageDB

	rows, err := db.DB.Query(`SELECT id, userId, message, isRead, createdAt FROM "message" WHERE chatId = $1`, chatId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var messageStruct types.MessageDB
		var id int
		var userId int
		var message string
		var isRead bool
		var createdAt time.Time

		rows.Scan(&id, &userId, &message, &isRead, &createdAt)

		messageStruct = types.MessageDB{
			Id:        id,
			UserId:    userId,
			Message:   message,
			IsRead:    isRead,
			CreatedAt: createdAt,
		}

		messages = append(messages, messageStruct)
	}

	return messages
}

func (db *DB) MarkAsRead(messageId int) {
	db.DB.Exec(`UPDATE "message" SET isread = true WHERE id = $1`, messageId)
}

func (db *DB) UserExist(userA, userB int) bool {
	var exist bool
	row := db.DB.QueryRow(`SELECT COUNT(*) = 2 AS both_exist FROM "user" WHERE id IN ($1, $2);`, userA, userB)

	row.Scan(&exist)

	return exist
}

func (db *DB) IsOnline(user int) bool {
	var isOnline bool

	row := db.DB.QueryRow(`SELECT isonline FROM "user" WHERE id=$1`, user)

	row.Scan(&isOnline)

	return isOnline
}

func (db *DB) SaveMessage(chatId int, sender int, message string) {
	_, err := db.DB.Exec(`INSERT INTO "message" (userid, chatid, message) VALUES ($1, $2, $3)`, sender, chatId, message)

	if err != nil {
		panic(err)
	}
}
