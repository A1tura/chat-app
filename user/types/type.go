package types

import "github.com/dgrijalva/jwt-go"


type RegisterReq struct {
    Username string `json:"username"`
	Password string `json:"password"`
}

type AuthClaims struct {
    Id string `json:"id"`
    jwt.StandardClaims
}
