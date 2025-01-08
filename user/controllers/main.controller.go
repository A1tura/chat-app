package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	db "user/dal"
	"user/types"

	"github.com/dgrijalva/jwt-go"
)

func RegisterController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		db := r.Context().Value("db").(*db.DB)
		var user types.RegisterReq
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			panic(err)
		}

		if user.Password == "" || user.Username == "" {
			fmt.Fprintf(w, "You must specify all the fields")
			return
		}

		if db.UserExist(user.Username) == true {
			fmt.Fprintf(w, "Username is already taken!")
			return
		}

		password_hash_bytes := sha512.New()
		password_hash_bytes.Write([]byte(user.Password))
		password_hash := hex.EncodeToString(password_hash_bytes.Sum(nil))

		created, id := db.CreateUser(user.Username, password_hash)

		if created {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id": strconv.Itoa(id),
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				panic(err)
			}
			w.Header().Add("Authorization", tokenString)
			fmt.Fprintf(w, "Registered successful")
		} else {
			fmt.Fprintf(w, "Error")
		}

	} else {
		http.NotFound(w, r)
	}
}

func SigninController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		defer r.Body.Close()
		db := r.Context().Value("db").(*db.DB)
		var user types.RegisterReq
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			panic(err)
		}

		if user.Password == "" || user.Username == "" {
			fmt.Fprintf(w, "You must specify all the fields")
			return
		}

		password_hash_bytes := sha512.New()
		password_hash_bytes.Write([]byte(user.Password))
		password_hash := hex.EncodeToString(password_hash_bytes.Sum(nil))

		exist, id := db.SigninUser(user.Username, password_hash)

		if exist == false {
			http.Error(w, "User with this creds do not exist", 403)
			return
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id": strconv.Itoa(id),
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			if err != nil {
				panic(err)
			}
			w.Header().Add("Authorization", tokenString)
			fmt.Fprint(w, "Successful login")
		}

	} else {
		http.NotFound(w, r)
	}
}
