package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func Connect(con string) *DB {
	db, err := sql.Open("postgres", con)
	if err != nil {
		panic(err)
	}

	if db.Ping(); err != nil {
		panic(err)
	}

	return &DB{DB: db}
}

func (db *DB) UserExist(username string) bool {

	row := db.QueryRow(`SELECT 1 FROM "user" WHERE username = $1;`, username)

	var exist int
	row.Scan(&exist)

	if exist == 1 {
		return true
	}

	return false
}

func (db *DB) CreateUser(username string, password_hash string) (bool, int) {

	row := db.QueryRow(`INSERT INTO "user" (username, password_hash, createdat) VALUES ($1, $2, $3) RETURNING id;`, username, password_hash, time.DateTime)

	var id int
	row.Scan(&id)

	return true, id
}

func (db *DB) SigninUser(username string, password_hash string) (bool, int) {
	row := db.QueryRow(`SELECT id FROM "user" WHERE username = $1 AND password_hash = $2`, username, password_hash)

	var id int
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, -1
		}
	}

	return true, id
}
